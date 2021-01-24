package ssinterface

import (
	"context"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type CycleJobService interface {
	DoJob(ctx context.Context, logger interfaces.Logger) (time.Duration, error)
}
