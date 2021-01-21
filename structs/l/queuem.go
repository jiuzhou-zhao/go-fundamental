package l

import "container/list"

type QueueM struct {
	l *list.List
	m map[interface{}]*list.Element
}

func NewQueueM() *QueueM {
	return &QueueM{
		l: list.New(),
		m: make(map[interface{}]*list.Element),
	}
}

type qi struct {
	k interface{}
	v interface{}
}

func (qm *QueueM) tryRemoveByKey(k interface{}) interface{} {
	if e, ok := qm.m[k]; ok {
		qm.l.Remove(e)
		delete(qm.m, k)
		return e.Value.(*qi).v
	}
	return nil
}

func (qm *QueueM) Size() int {
	return qm.l.Len()
}

func (qm *QueueM) Update(k, v interface{}) bool {
	if e, ok := qm.m[k]; ok {
		e.Value.(*qi).v = v
		return true
	}
	return false
}

func (qm *QueueM) PushBack(k, v interface{}) {
	qm.tryRemoveByKey(k)
	qm.m[k] = qm.l.PushBack(&qi{
		k: k,
		v: v,
	})
}

func (qm *QueueM) PushFront(k, v interface{}) {
	qm.tryRemoveByKey(k)
	qm.m[k] = qm.l.PushFront(&qi{
		k: k,
		v: v,
	})
}

func (qm *QueueM) PopBack() (k, v interface{}) {
	e := qm.l.Back()
	if e != nil {
		k = e.Value.(*qi).k
		v = e.Value.(*qi).v
		qm.tryRemoveByKey(k)
	}
	return
}

func (qm *QueueM) PopFront() (k, v interface{}) {
	e := qm.l.Front()
	if e != nil {
		k = e.Value.(*qi).k
		v = e.Value.(*qi).v
		qm.tryRemoveByKey(k)
	}
	return
}

func (qm *QueueM) Front() (k, v interface{}) {
	e := qm.l.Front()
	if e != nil {
		k = e.Value.(*qi).k
		v = e.Value.(*qi).v
	}
	return
}

func (qm *QueueM) Back() (k, v interface{}) {
	e := qm.l.Back()
	if e != nil {
		k = e.Value.(*qi).k
		v = e.Value.(*qi).v
	}
	return
}

func (qm *QueueM) KeyExists(k interface{}) bool {
	_, ok := qm.m[k]
	return ok
}

func (qm *QueueM) GetValueByKey(k interface{}) interface{} {
	if e, ok := qm.m[k]; ok {
		return e.Value.(*qi).v
	}
	return nil
}

func (qm *QueueM) Remove(k interface{}) interface{} {
	return qm.tryRemoveByKey(k)
}

type IteratorSet struct {
	K interface{}
	V interface{}
}

func (qm *QueueM) Iterator() <-chan *IteratorSet {
	ch := make(chan *IteratorSet, 10)
	go func() {
		for e := qm.l.Front(); e != nil; e = e.Next() {
			ch <- &IteratorSet{
				K: e.Value.(*qi).k,
				V: e.Value.(*qi).v,
			}
		}
		close(ch)
	}()
	return ch
}
