package loge

import (
	"context"
	"log"
	"sync"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

var (
	globalLock   sync.RWMutex
	globalLogger *Logger
)

type Logger struct {
	loggerImpl interfaces.Logger
}

func NewLogger(logger interfaces.Logger) *Logger {
	if logger == nil {
		logger = &ConsoleLogger{}
	}
	return &Logger{loggerImpl: logger}
}

func (logger *Logger) GetLogger() interfaces.Logger {
	return logger.loggerImpl
}

func (logger *Logger) Debug(ctx context.Context, v ...interface{}) {
	logger.loggerImpl.Record(ctx, interfaces.LogLevelDebug, v...)
}

func (logger *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	logger.loggerImpl.Recordf(ctx, interfaces.LogLevelDebug, format, v...)
}

func (logger *Logger) Info(ctx context.Context, v ...interface{}) {
	logger.loggerImpl.Record(ctx, interfaces.LogLevelInfo, v...)
}

func (logger *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	logger.loggerImpl.Recordf(ctx, interfaces.LogLevelInfo, format, v...)
}

func (logger *Logger) Warn(ctx context.Context, v ...interface{}) {
	logger.loggerImpl.Record(ctx, interfaces.LogLevelWarn, v...)
}

func (logger *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	logger.loggerImpl.Recordf(ctx, interfaces.LogLevelWarn, format, v...)
}

func (logger *Logger) Error(ctx context.Context, v ...interface{}) {
	logger.loggerImpl.Record(ctx, interfaces.LogLevelError, v...)
}

func (logger *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	logger.loggerImpl.Recordf(ctx, interfaces.LogLevelError, format, v...)
}

func (logger *Logger) Fatal(ctx context.Context, v ...interface{}) {
	logger.loggerImpl.Record(ctx, interfaces.LogLevelFatal, v...)
}

func (logger *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	logger.loggerImpl.Recordf(ctx, interfaces.LogLevelFatal, format, v...)
}

func SetGlobalLogger(logger *Logger) *Logger {
	globalLock.Lock()
	defer globalLock.Unlock()

	oldLogger := globalLogger
	globalLogger = logger
	return oldLogger
}

func GetGlobalLogger() *Logger {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return globalLogger
}

func _logOrPanic(f func(logger *Logger)) {
	logger := GetGlobalLogger()
	if logger == nil {
		log.Fatal("no global logger")
	}
	f(logger)
}

func Debug(ctx context.Context, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Debug(ctx, v...)
	})
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Debugf(ctx, format, v...)
	})
}

func Info(ctx context.Context, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Info(ctx, v...)
	})
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Infof(ctx, format, v...)
	})
}

func Warn(ctx context.Context, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Warn(ctx, v...)
	})
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Warnf(ctx, format, v...)
	})
}

func Error(ctx context.Context, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Error(ctx, v...)
	})
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Errorf(ctx, format, v...)
	})
}

func Fatal(ctx context.Context, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Fatal(ctx, v...)
	})
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	_logOrPanic(func(logger *Logger) {
		logger.Fatalf(ctx, format, v...)
	})
}
