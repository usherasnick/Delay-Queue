package delayqueue

import (
	"context"

	"github.com/rs/zerolog/log"
)

type RedisRWRequest struct {
	RequestType RedisRequestType
	RequestOp   RedisRWRequestOp
	Inputs      []interface{}
	ResponseCh  chan *RedisRWResponse
}

type RedisRWResponse struct {
	Outputs []interface{}
	Err     error
}

type RedisRequestType int

const (
	TopicRequest  RedisRequestType = 1
	BucketRequest RedisRequestType = 2
)

type RedisRWRequestOp int

const (
	PutTopicRequest  RedisRWRequestOp = 1
	ListTopicRequest RedisRWRequestOp = 2
	HasTopicRequest  RedisRWRequestOp = 3
	DelTopicRequest  RedisRWRequestOp = 4

	PushToBucketRequest     RedisRWRequestOp = 5
	GetOneFromBucketRequest RedisRWRequestOp = 6
	DelFromBucketRequest    RedisRWRequestOp = 7
)

// 利用命令管道来提高服务的响应速度, 并且提供更高层面的数据一致性.
func (dq *DelayQueue) sendRedisRWRequest(req *RedisRWRequest) {
	switch req.RequestType {
	case TopicRequest:
		{
			dq.topicRWChannel <- req
		}
	case BucketRequest:
		{
			dq.bucketRWChannel <- req
		}
	default:
		{
			log.Error().Msgf("invalid RedisRequestType")
		}
	}
}

func (dq *DelayQueue) handleTopicRWRequest(ctx context.Context) {
REDIS_RW_LOOP:
	for {
		select {
		case <-ctx.Done():
			{
				break REDIS_RW_LOOP
			}
		case req := <-dq.topicRWChannel:
			{
				if req.RequestOp == PutTopicRequest {
					err := dq.putTopic(req.Inputs[0].(string), req.Inputs[1].(string))
					req.ResponseCh <- &RedisRWResponse{
						Err: err,
					}
				} else if req.RequestOp == ListTopicRequest {
					v, err := dq.listTopic(req.Inputs[0].(string))
					if err != nil {
						req.ResponseCh <- &RedisRWResponse{
							Err: err,
						}
					} else {
						req.ResponseCh <- &RedisRWResponse{
							Outputs: []interface{}{v},
							Err:     nil,
						}
					}
				} else if req.RequestOp == HasTopicRequest {
					v, err := dq.hasTopic(req.Inputs[0].(string), req.Inputs[1].(string))
					if err != nil {
						req.ResponseCh <- &RedisRWResponse{
							Err: err,
						}
					} else {
						req.ResponseCh <- &RedisRWResponse{
							Outputs: []interface{}{v},
							Err:     nil,
						}
					}
				} else if req.RequestOp == DelTopicRequest {
					err := dq.delTopic(req.Inputs[0].(string), req.Inputs[1].(string))
					req.ResponseCh <- &RedisRWResponse{
						Err: err,
					}
				}
			}
		}
	}
}

func (dq *DelayQueue) handleBucketRWRequest(ctx context.Context) {
REDIS_RW_LOOP:
	for {
		select {
		case <-ctx.Done():
			{
				break REDIS_RW_LOOP
			}
		case req := <-dq.bucketRWChannel:
			{
				if req.RequestOp == PushToBucketRequest {
					err := dq.pushToBucket(req.Inputs[0].(string), req.Inputs[1].(int64), req.Inputs[2].(string))
					req.ResponseCh <- &RedisRWResponse{
						Err: err,
					}
				} else if req.RequestOp == GetOneFromBucketRequest {
					v, err := dq.getOneFromBucket(req.Inputs[0].(string))
					if err != nil {
						req.ResponseCh <- &RedisRWResponse{
							Err: err,
						}
					} else {
						req.ResponseCh <- &RedisRWResponse{
							Outputs: []interface{}{v},
							Err:     nil,
						}
					}
				} else if req.RequestOp == DelFromBucketRequest {
					err := dq.delFromBucket(req.Inputs[0].(string), req.Inputs[1].(string))
					req.ResponseCh <- &RedisRWResponse{
						Err: err,
					}
				}
			}
		}
	}
}
