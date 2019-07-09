package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	conf "github.com/usherasnick/Delay-Queue/internal/config"
	"github.com/usherasnick/Delay-Queue/internal/util"
)

var (
	cfgPathFlag = flag.String("conf", "conf/delay_queue.json", "delay queue config file")
	verboseFlag = flag.Bool("verbose", false, "set verbose output")
)

func main() {
	flag.Parse()

	// 设置全局logger
	logOutput := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logOutput.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf("[%v]", i)
	}
	logOutput.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%-6s]", i))
	}
	logOutput.FormatFieldName = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("<%s: ", i))
	}
	logOutput.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s>", i)
	}
	logOutput.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	log.Logger = log.Output(logOutput)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verboseFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// 加载节点配置
	var cfg conf.DelayQueueService
	util.LoadConfigFileOrPanic(*cfgPathFlag, &cfg)

	// 生成全局上下文控制器
	gctx, gcancel := context.WithCancel(context.Background())

	srv := NewTaskDelayQueueServiceServer(&cfg)
	defer func() {
		gcancel()
		srv.Close()
	}()

	// 开启grpc服务
	go serveGPRC(gctx, srv, cfg.GRPCEndpoint)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	for range sigCh {
		break
	}

	log.Info().Msg("task-delay-queue-service has been stopped")
}
