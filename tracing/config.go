package tracing

import (
	"time"
)

// Config is used for Tracer creation
type Config struct {
	ServerAddr    string        `yaml:"server_addr"`
	ServiceName   string        `yaml:"service_name"`
	FlushInterval time.Duration `yaml:"flush_interval"`
}
