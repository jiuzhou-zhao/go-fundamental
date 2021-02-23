package tools

import "github.com/jiuzhou-zhao/go-fundamental/structs/l"

type LruList interface {
	Add(k, v interface{})
	Remove(k interface{})
	Exists(k interface{}) bool
}

type lruList struct {
	queue   *l.QueueM
	maxSize int
}

func NewLruList(maxSize int) LruList {
	return &lruList{
		queue:   l.NewQueueM(),
		maxSize: maxSize,
	}
}

func (l *lruList) Add(k, v interface{}) {
	if l.maxSize > 0 && l.queue.Size() >= l.maxSize {
		l.queue.PopFront()
	}
	l.queue.PushBack(k, v)
}

func (l *lruList) Remove(k interface{}) {
	l.queue.Remove(k)
}

func (l *lruList) Exists(k interface{}) bool {
	return l.queue.KeyExists(k)
}
