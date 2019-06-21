package delayqueue

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	pb "github.com/usherasnick/Delay-Queue/api"
	conf "github.com/usherasnick/Delay-Queue/internal/config"
	"github.com/usherasnick/Delay-Queue/internal/kafka"
	"github.com/usherasnick/Delay-Queue/internal/redis"
)

type DelayQueue struct {
	ctx    context.Context
	cancel context.CancelFunc

	topicRWChannel  chan *RedisRWRequest
	bucketRWChannel chan *RedisRWRequest
	readyQ          *ReadyQueue

	redisCli *redis.RedisConnPoolSingleton
	producer *kafka.Producer
}

func NewDelayQueue(cfg *conf.DelayQueueService) *DelayQueue {
	ctx, cancel := context.WithCancel(context.Background())
	dq := &DelayQueue{
		ctx:    ctx,
		cancel: cancel,

		topicRWChannel:  make(chan *RedisRWRequest, 1024),
		bucketRWChannel: make(chan *RedisRWRequest, 1024),
		readyQ:          NewReadyQueue(),

		redisCli: redis.GetOrCreateInstance(cfg.RedisService),
		producer: kafka.NewProducer(cfg.KafkaService),
	}
	go dq.handleTopicRWRequest(ctx)
	go dq.handleBucketRWRequest(ctx)
	go dq.handleTimer(ctx)
	go dq.poll(ctx)
	return dq
}

func (dq *DelayQueue) Close() {
	redis.ReleaseInstance()
	if dq.cancel != nil {
		dq.cancel()
	}
	dq.producer.Close()
}

// DelayQueue的Push/PushTopic/RemoveTopic按照原来的处理流程, 可能会受主循环影响而被阻塞
// 因此作乐观处理, 不等待处理结果直接返回, 由调用方确认是否操作成功

func (dq *DelayQueue) Push(task *Task) error {
	/* start to add task */
	err := dq.putTask(task.Id, task)
	if err != nil {
		log.Error().Err(err).Msgf("failed to add task <id: %s>", task.Id)
		return err
	}

	/* start to push task into bucket */
	resp := make(chan *RedisRWResponse, 1)
	req := &RedisRWRequest{
		RequestType: BucketRequest,
		RequestOp:   PushToBucketRequest,
		Inputs:      []interface{}{DefaultBucketName, task.Delay, task.Id},
		ResponseCh:  resp,
	}
	dq.sendRedisRWRequest(req)
	log.Debug().Msgf("add a new task <%s>", task.Id)
	return nil
}

func (dq *DelayQueue) Remove(taskId string) error {
	/* start to delete task */
	err := dq.delTask(taskId)
	if err != nil {
		log.Error().Err(err).Msgf("failed to remove task <id: %s>", taskId)
		return err
	}
	log.Debug().Msgf("delete a task <%s>", taskId)
	return nil
}

func (dq *DelayQueue) Get(taskId string) (*Task, error) {
	/* start to get task */
	task, err := dq.getTask(taskId)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get task <id: %s>", taskId)
		return nil, err
	}
	// 任务不存在, 可能已被删除
	if task == nil {
		return nil, nil
	}
	return task, nil
}

func (dq *DelayQueue) PushTopic(topic string) error {
	/* start to add topic */
	resp := make(chan *RedisRWResponse, 1)
	req := &RedisRWRequest{
		RequestType: TopicRequest,
		RequestOp:   PutTopicRequest,
		Inputs:      []interface{}{DefaultTopicSetName, topic},
		ResponseCh:  resp,
	}
	dq.sendRedisRWRequest(req)
	return nil
}

func (dq *DelayQueue) RemoveTopic(topic string) error {
	/* start to delete topic */
	resp := make(chan *RedisRWResponse, 1)
	req := &RedisRWRequest{
		RequestType: TopicRequest,
		RequestOp:   DelTopicRequest,
		Inputs:      []interface{}{DefaultTopicSetName, topic},
		ResponseCh:  resp,
	}
	dq.sendRedisRWRequest(req)
	return nil
}

func (dq *DelayQueue) handleTimer(ctx context.Context) {
	timer := time.NewTicker(1 * time.Second)
TICK_LOOP:
	for {
		select {
		case <-ctx.Done():
			{
				break TICK_LOOP
			}
		case t := <-timer.C:
			{
				dq.timerHandler(t.Unix())
			}
		}
	}
}

func (dq *DelayQueue) timerHandler(now int64) {
	for {
		/* start to get task from bucket */
		resp := make(chan *RedisRWResponse)
		req := &RedisRWRequest{
			RequestType: BucketRequest,
			RequestOp:   GetOneFromBucketRequest,
			Inputs:      []interface{}{DefaultBucketName},
			ResponseCh:  resp,
		}
		dq.sendRedisRWRequest(req)
		outs := <-resp
		if outs.Err != nil {
			log.Error().Err(outs.Err).Msgf("failed to scan bucket <name: %s>", DefaultBucketName)
			return
		}
		bucketItem := outs.Outputs[0].(*BucketItem)
		if bucketItem == nil {
			return
		}

		// 延迟执行时间未到
		if bucketItem.TaskTimestamp > now {
			return
		}

		// 延迟执行时间小于等于当前时间, 取出任务并放入ReadyQueue
		/* start to get task */
		task, err := dq.getTask(bucketItem.TaskId)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get task <id: %s>", bucketItem.TaskId)
			continue
		}
		// 任务不存在, 可能已被删除, 马上从bucket中删除
		if task == nil {
			/* start to delete task from bucket */
			resp := make(chan *RedisRWResponse, 1)
			req := &RedisRWRequest{
				RequestType: BucketRequest,
				RequestOp:   DelFromBucketRequest,
				Inputs:      []interface{}{DefaultBucketName, bucketItem.TaskId},
				ResponseCh:  resp,
			}
			dq.sendRedisRWRequest(req)
			continue
		}
		// 再次确认任务延迟执行时间是否小于等于当前时间
		if task.Delay <= now {
			/* start to push task into ready queue */
			dq.readyQ.PushToReadyQueue(dq.redisCli, DefaultReadyQueueName, task.Id) // nolint

			/* start to delete task from bucket */
			resp = make(chan *RedisRWResponse, 1)
			req = &RedisRWRequest{
				RequestType: BucketRequest,
				RequestOp:   DelFromBucketRequest,
				Inputs:      []interface{}{DefaultBucketName, task.Id},
				ResponseCh:  resp,
			}
			dq.sendRedisRWRequest(req)
		}
	}
}

func (dq *DelayQueue) poll(ctx context.Context) {
POLL_LOOP:
	for {
		select {
		case <-ctx.Done():
			{
				break POLL_LOOP
			}
		default:
			{
				/* start to get all subscribed topics */
				resp := make(chan *RedisRWResponse)
				req := &RedisRWRequest{
					RequestType: TopicRequest,
					RequestOp:   ListTopicRequest,
					Inputs:      []interface{}{DefaultTopicSetName},
					ResponseCh:  resp,
				}
				dq.sendRedisRWRequest(req)
				outs := <-resp
				if outs.Err != nil {
					log.Error().Err(outs.Err).Msg("failed to list topic")
					continue
				}
				topics := outs.Outputs[0].([]string)
				if len(topics) == 0 {
					continue
				}
				log.Debug().Msgf("all subscribed topics: %v", topics)

				/* start to pop task from ready queue */
				// TODO: 在无任务和有任务两种状态之间切换会带来额外的延时, 可能会影响具体的业务
				taskId, err := dq.readyQ.BlockPopFromReadyQueue(dq.redisCli, DefaultReadyQueueName, DefaultBlockPopFromReadyQueueTimeout)
				if err != nil {
					log.Error().Err(err).Msg("failed to pop from ready queue")
					continue
				}
				if taskId == "" {
					continue
				}
				log.Debug().Msgf("get ready task <%s>", taskId)

				/* start to get task */
				task, err := dq.getTask(taskId)
				if err != nil {
					log.Error().Err(err).Msgf("failed to get task <id: %s>", taskId)
					continue
				}
				// 任务不存在, 可能已被删除
				if task == nil {
					continue
				}

				/* check whether the task has been subscribed or not, if not, just pass it */
				subscribed := false
				for _, topic := range topics {
					if topic == task.Topic {
						subscribed = true
						break
					}
				}
				if !subscribed {
					continue
				}

				// TTR的设计目的是为了保证消息传输的可靠性
				// 任务执行完成后, 消费端需要调用finish接口去删除任务, 否则任务会重复投递, 消费端必须能处理同一任务多次投递的情形
				timestamp := time.Now().Unix() + task.TTR
				/* start to push task into bucket */
				resp = make(chan *RedisRWResponse)
				req = &RedisRWRequest{
					RequestType: BucketRequest,
					RequestOp:   PushToBucketRequest,
					Inputs:      []interface{}{DefaultBucketName, timestamp, task.Id},
					ResponseCh:  resp,
				}
				dq.sendRedisRWRequest(req)

				/* start to publish ready task to kafka */
				msg, _ := proto.Marshal(&pb.Task{
					Id:            task.Id,
					AttachedTopic: task.Topic,
					Payload:       task.Blob,
				})
				if err := dq.producer.Publish(task.Topic, msg); err != nil {
					log.Error().Err(err).Msgf("failed to publish ready task <id: %s>", task.Id)
				}
				log.Debug().Msgf("send ready task <%s> to kafka", taskId)
			}
		}
	}
}
