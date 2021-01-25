package interfaces

import (
	"context"
)

type LoggerLevel int

const (
	LogLevelDebug LoggerLevel = 0x01
	LogLevelInfo  LoggerLevel = 0x02
	LogLevelWarn  LoggerLevel = 0x04
	LogLevelError LoggerLevel = 0x08
	LogLevelFatal LoggerLevel = 0x10
)

func (ll LoggerLevel) String() string {
	switch ll {
	case LogLevelDebug:
		return "DBG"
	case LogLevelInfo:
		return "INF"
	case LogLevelWarn:
		return "WRN"
	case LogLevelError:
		return "ERR"
	case LogLevelFatal:
		return "FTL"
	}
	return "UKN"
}

type Logger interface {
	Record(ctx context.Context, level LoggerLevel, v ...interface{})
	Recordf(ctx context.Context, level LoggerLevel, format string, v ...interface{})
}
