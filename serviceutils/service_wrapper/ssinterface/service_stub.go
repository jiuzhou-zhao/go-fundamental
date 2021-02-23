package ssinterface

import (
	"context"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type ServiceStub interface {
	Run(ctx context.Context, logger interfaces.Logger)
}
