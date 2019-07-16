package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/usherasnick/Delay-Queue/api"
)

func PushTaskPreCheck(req *pb.PushTaskRequest) error {
	if req.GetTask().GetId() == "" {
		return status.Errorf(codes.InvalidArgument, "empty task id")
	}
	if req.GetTask().GetAttachedTopic() == "" {
		return status.Errorf(codes.InvalidArgument, "empty attached topic")
	}
	if req.GetTask().GetDelay() <= 0 {
		return status.Errorf(codes.InvalidArgument, "invalid task delay")
	}
	if req.GetTask().GetTtr() <= 0 || req.GetTask().GetTtr() > 86400 {
		return status.Errorf(codes.InvalidArgument, "empty attached topic")
	}
	return nil
}

func FinishTaskPreCheck(req *pb.FinishTaskRequest) error {
	if req.GetTaskId() == "" {
		return status.Errorf(codes.InvalidArgument, "empty task id")
	}
	return nil
}

func CheckTaskPreCheck(req *pb.CheckTaskRequest) error {
	if req.GetTaskId() == "" {
		return status.Errorf(codes.InvalidArgument, "empty task id")
	}
	return nil
}

func SubscribeTopicPreCheck(req *pb.SubscribeTopicRequest) error {
	if req.GetTopic() == "" {
		return status.Errorf(codes.InvalidArgument, "empty topic")
	}
	return nil
}

func UnsubscribeTopicPreCheck(req *pb.UnsubscribeTopicRequest) error {
	if req.GetTopic() == "" {
		return status.Errorf(codes.InvalidArgument, "empty topic")
	}
	return nil
}

func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
