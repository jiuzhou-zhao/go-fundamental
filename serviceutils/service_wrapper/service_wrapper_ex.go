package service_wrapper

import (
	"context"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/serviceutils/service_wrapper/ssinterface"
	"github.com/jiuzhou-zhao/go-fundamental/serviceutils/service_wrapper/stub"
)

type CycleServiceWrapper struct {
	*ServiceWrapper
}

func NewCycleServiceWrapper(ctx context.Context, logger interfaces.Logger) *CycleServiceWrapper {
	return &CycleServiceWrapper{
		ServiceWrapper: NewServiceWrapper(ctx, logger),
	}
}

func (sw *CycleServiceWrapper) Start(serviceImpl ssinterface.CycleJobService) error {
	return sw.ServiceWrapper.Start(stub.NewCycleJobServiceStub(serviceImpl))
}
