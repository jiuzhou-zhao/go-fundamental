package servicetoolset

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
)

func SignalContext(ctx context.Context, logger interfaces.Logger) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		eLog := loge.NewLogger(logger)
		eLog.Info(ctx, "listening for shutdown signal")
		<-sigs
		eLog.Info(ctx, "shutdown signal received")
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}

type AbstractServer interface {
	Run(ctx context.Context) error
}

type ServerHelper struct {
	ctx    context.Context
	wg     sync.WaitGroup
	logger *loge.Logger
}

func NewServerHelper(ctx context.Context, logger interfaces.Logger) *ServerHelper {
	return &ServerHelper{
		ctx:    SignalContext(ctx, logger),
		logger: loge.NewLogger(logger),
	}
}

func (sh *ServerHelper) StartServer(s AbstractServer) {
	sh.wg.Add(1)
	go func() {
		defer sh.wg.Done()
		if err := s.Run(sh.ctx); err != nil {
			sh.logger.Fatalf(context.Background(), "runServer error:%v", err)
		}
	}()
}

func (sh *ServerHelper) Wait() {
	sh.wg.Wait()
}
