package tools

import (
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/structs/l"
)

type ExistsChecker interface {
	Add(key interface{})
	Exists(key interface{}) bool
}

type ExistsCheckerWithMaxSize struct {
	queue   *l.QueueM
	maxSize int
}

func NewExistsCheckerWithMaxSize(maxSize int) ExistsChecker {
	return &ExistsCheckerWithMaxSize{
		queue:   l.NewQueueM(),
		maxSize: maxSize,
	}
}

func (ec *ExistsCheckerWithMaxSize) Add(key interface{}) {
	if ec.maxSize > 0 && ec.queue.Size() >= ec.maxSize {
		ec.queue.PopFront()
	}
	ec.queue.PushBack(key, time.Now())
}

func (ec *ExistsCheckerWithMaxSize) Exists(key interface{}) bool {
	return ec.queue.KeyExists(key)
}
