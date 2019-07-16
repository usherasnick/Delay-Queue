package main

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/usherasnick/Delay-Queue/api"
	conf "github.com/usherasnick/Delay-Queue/internal/config"
	delayqueue "github.com/usherasnick/Delay-Queue/internal/delay-queue"
)

type TaskDelayQueueServiceServer struct {
	pb.UnimplementedTaskDelayQueueServiceServer
	inst *delayqueue.DelayQueue
}

func NewTaskDelayQueueServiceServer(cfg *conf.DelayQueueService) *TaskDelayQueueServiceServer {
	return &TaskDelayQueueServiceServer{
		inst: delayqueue.NewDelayQueue(cfg),
	}
}

func (srv *TaskDelayQueueServiceServer) Close() {
	srv.inst.Close()
}

func (srv *TaskDelayQueueServiceServer) PushTask(ctx context.Context, req *pb.PushTaskRequest) (*pb.PushTaskResponse, error) {
	if err := PushTaskPreCheck(req); err != nil {
		return nil, err
	}

	task := &delayqueue.Task{}
	task.Id = req.GetTask().GetId()
	task.Topic = req.GetTask().GetAttachedTopic()
	task.Delay = time.Now().Unix() + int64(req.GetTask().GetDelay()) // 当前时间戳 + 相对时间 = 绝对时间
	task.TTR = int64(req.GetTask().GetTtr())
	task.Blob = req.GetTask().GetPayload()

	if err := srv.inst.Push(task); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if IsContextDone(ctx) {
		return nil, status.Errorf(codes.DeadlineExceeded, "system may be under busy")
	}

	return &pb.PushTaskResponse{}, nil
}

func (srv *TaskDelayQueueServiceServer) FinishTask(ctx context.Context, req *pb.FinishTaskRequest) (*pb.FinishTaskResponse, error) {
	if err := FinishTaskPreCheck(req); err != nil {
		return nil, err
	}

	if err := srv.inst.Remove(req.GetTaskId()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if IsContextDone(ctx) {
		return nil, status.Errorf(codes.DeadlineExceeded, "system may be under busy")
	}

	return &pb.FinishTaskResponse{}, nil
}

func (srv *TaskDelayQueueServiceServer) CheckTask(ctx context.Context, req *pb.CheckTaskRequest) (*pb.CheckTaskResponse, error) {
	if err := CheckTaskPreCheck(req); err != nil {
		return nil, err
	}

	task, err := srv.inst.Get(req.GetTaskId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if IsContextDone(ctx) {
		return nil, status.Errorf(codes.DeadlineExceeded, "system may be under busy")
	}
	if task == nil {
		return &pb.CheckTaskResponse{}, nil
	}

	return &pb.CheckTaskResponse{
		Task: &pb.Task{
			Id:            task.Id,
			AttachedTopic: task.Topic,
			Ttr:           int32(task.TTR),
			Payload:       task.Blob,
		},
	}, nil
}

func (srv *TaskDelayQueueServiceServer) SubscribeTopic(ctx context.Context, req *pb.SubscribeTopicRequest) (*pb.SubscribeTopicResponse, error) {
	if err := SubscribeTopicPreCheck(req); err != nil {
		return nil, err
	}

	if err := srv.inst.PushTopic(req.GetTopic()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if IsContextDone(ctx) {
		return nil, status.Errorf(codes.DeadlineExceeded, "system may be under busy")
	}

	return &pb.SubscribeTopicResponse{}, nil
}

func (srv *TaskDelayQueueServiceServer) UnsubscribeTopic(ctx context.Context, req *pb.UnsubscribeTopicRequest) (*pb.UnsubscribeTopicResponse, error) {
	if err := UnsubscribeTopicPreCheck(req); err != nil {
		return nil, err
	}

	if err := srv.inst.RemoveTopic(req.GetTopic()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if IsContextDone(ctx) {
		return nil, status.Errorf(codes.DeadlineExceeded, "system may be under busy")
	}

	return &pb.UnsubscribeTopicResponse{}, nil
}
