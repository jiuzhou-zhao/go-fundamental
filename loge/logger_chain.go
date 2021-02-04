package loge

import (
	"context"
	"sync"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type LoggerChain struct {
	sync.RWMutex
	loggers    []interfaces.Logger
	loggersMap map[interfaces.Logger]interface{}
}

func NewLoggerChain() *LoggerChain {
	return &LoggerChain{
		loggers:    make([]interfaces.Logger, 0, 6),
		loggersMap: make(map[interfaces.Logger]interface{}),
	}
}

func (logger *LoggerChain) AppendLogger(log interfaces.Logger) {
	logger.Lock()
	defer logger.Unlock()
	if _, ok := logger.loggersMap[log]; ok {
		return
	}
	logger.loggersMap[log] = time.Now()
	logger.loggers = append(logger.loggers, log)
}

func (logger *LoggerChain) Record(ctx context.Context, depth int, level interfaces.LoggerLevel, v ...interface{}) {
	logger.RLock()
	defer logger.RUnlock()
	for _, log := range logger.loggers {
		log.Record(ctx, depth+1, level, v...)
	}
}

func (logger *LoggerChain) Recordf(ctx context.Context, depth int, level interfaces.LoggerLevel, format string, v ...interface{}) {
	logger.RLock()
	defer logger.RUnlock()
	for _, log := range logger.loggers {
		log.Recordf(ctx, depth+1, level, format, v...)
	}
}
