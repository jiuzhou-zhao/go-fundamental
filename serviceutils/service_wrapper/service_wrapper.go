package service_wrapper

import (
	"context"
	"errors"
	"sync"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/serviceutils/service_wrapper/ssinterface"
)

type ServiceWrapper struct {
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    interfaces.Logger

	validFlag bool
}

func NewServiceWrapper(ctx context.Context, logger interfaces.Logger) *ServiceWrapper {
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = &interfaces.ConsoleLogger{}
	}
	return &ServiceWrapper{
		ctx:       ctx,
		ctxCancel: cancel,
		logger:    logger,
		validFlag: true,
	}
}

func (sw *ServiceWrapper) Start(serviceImpl ssinterface.ServiceStub) error {
	if !sw.validFlag {
		return errors.New("service wrapper stop or closed")
	}
	if serviceImpl == nil {
		return errors.New("invalid input parameters")
	}
	sw.wg.Add(1)
	go func() {
		defer sw.wg.Done()
		serviceImpl.Run(sw.ctx, sw.logger)
	}()
	return nil
}

func (sw *ServiceWrapper) Stop() {
	sw.validFlag = false
	sw.ctxCancel()
}

func (sw *ServiceWrapper) Wait() {
	sw.validFlag = false
	sw.wg.Wait()
}
