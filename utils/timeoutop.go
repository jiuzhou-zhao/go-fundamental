package utils

import (
	"context"
	"time"
)

func TimeoutOp(ctx context.Context, timeout time.Duration, cb func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cb(ctx)
}

func DefRedisTimeoutOp(cb func(ctx context.Context)) {
	DefRedisTimeoutOpEx(context.Background(), cb)
}

func DefRedisTimeoutOpEx(ctx context.Context, cb func(ctx context.Context)) {
	TimeoutOp(ctx, 10*time.Second, cb)
}
