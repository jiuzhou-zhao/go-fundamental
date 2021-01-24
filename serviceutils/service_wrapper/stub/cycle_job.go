package stub

import (
	"context"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
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
	logger.Record(ctx, interfaces.LogLevelDebug, "enter cycle job loop")

	duration, err := ss.serviceImpl.DoJob(ctx, logger)
	if err != nil {
		logger.Recordf(ctx, interfaces.LogLevelError, "do job failed: %v", err)
		return
	}
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			logger.Record(ctx, interfaces.LogLevelDebug, "check ctx done, try exit loop")
		case <-time.After(duration):
			duration, err = ss.serviceImpl.DoJob(ctx, logger)
			if err != nil {
				logger.Recordf(ctx, interfaces.LogLevelError, "do job failed: %v", err)
				loop = false
				break
			}
		}
	}
	logger.Record(ctx, interfaces.LogLevelDebug, "leave cycle job loop")
}
