package tracing

import (
	"time"
)

// Config is used for Tracer creation
type Config struct {
	ServerAddr    string
	ServiceName   string
	FlushInterval time.Duration
}
