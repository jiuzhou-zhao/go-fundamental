package loge

import (
	"context"
	"fmt"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
)

type EmptyLogger struct {
}

func (logger *EmptyLogger) Record(ctx context.Context, level interfaces.LoggerLevel, v ...interface{}) {

}

func (logger *EmptyLogger) Recordf(ctx context.Context, level interfaces.LoggerLevel, format string, v ...interface{}) {

}

type ConsoleLogger struct {
}

func (l *ConsoleLogger) Record(ctx context.Context, level interfaces.LoggerLevel, v ...interface{}) {
	i := fmt.Sprint(v...)
	fmt.Printf("[%v] %v\n", level.String(), i)
	if level == interfaces.LogLevelFatal {
		panic(i)
	}
}

func (l *ConsoleLogger) Recordf(ctx context.Context, level interfaces.LoggerLevel, format string, v ...interface{}) {
	i := fmt.Sprintf(format, v...)
	fmt.Printf("[%v] %v\n", level.String(), i)
	if level == interfaces.LogLevelFatal {
		panic(i)
	}
}
