package main

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/usherasnick/Delay-Queue/api"
)

func serveGPRC(ctx context.Context, srv *TaskDelayQueueServiceServer, ep string) {
	l, err := net.Listen("tcp", ep)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task-delay-queue-grpc-service")
	}

	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(8 * 1024 * 1024), // 限定单次请求最大可发送8Mb数据
		grpc.MaxRecvMsgSize(8 * 1024 * 1024), // 限定单次请求最大可接收8Mb数据
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true, // 允许非活跃流连接发送探活请求
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    5 * time.Minute,
			Timeout: 1 * time.Minute, // 允许已标记为非活跃客户端连接最长可存活的时间
		}),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTaskDelayQueueServiceServer(grpcServer, srv)
	// for mock testing with grpcurl
	reflection.Register(grpcServer)

	log.Info().Msgf("task-delay-queue-grpc-service is listening at \x1b[1;31m%s\x1b[0m", ep)
	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Warn().Err(err).Msgf("task-delay-queue-grpc-service will be terminated")
		}
	}()

GRPC_LOOP:
	for { // nolint
		select {
		case <-ctx.Done():
			break GRPC_LOOP
		}
	}

	grpcServer.GracefulStop()
	log.Info().Msg("task-delay-queue-grpc-service has been stopped")
}
