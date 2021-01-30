package stub

import (
	"context"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/jiuzhou-zhao/go-fundamental/serviceutils/service_wrapper/ssinterface"
)

type cycleJobService struct {
	serviceImpl ssinterface.CycleJobService
}

func NewCycleJobServiceStub(serviceImpl ssinterface.CycleJobService) ssinterface.ServiceStub {
	if serviceImpl == nil {
		return nil
	}
	return &cycleJobService{
		serviceImpl: serviceImpl,
	}
}

func (ss *cycleJobService) Run(ctx context.Context, logger interfaces.Logger) {
	loop := true

	eLog := loge.NewLogger(logger)
	eLog.Debug(ctx, "enter cycle job loop")

	duration, err := ss.serviceImpl.DoJob(ctx, logger)
	if err != nil {
		eLog.Errorf(ctx, "do job failed: %v", err)
		return
	}
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			eLog.Debug(ctx, "check ctx done, try exit loop")
		case <-time.After(duration):
			duration, err = ss.serviceImpl.DoJob(ctx, logger)
			if err != nil {
				eLog.Errorf(ctx, "do job failed: %v", err)
				loop = false
				break
			}
		}
	}
	eLog.Debug(ctx, "leave cycle job loop")
}
