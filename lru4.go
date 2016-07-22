// Author:dup2X
// Last modified: 2016-07-22 11:48
// Filename: lru4.go
package lru4

import (
	"container/list"
	"sync"
)

const defaultStep = 4

type LRU4Key interface{}

type LRU4Cache struct {
	MaxEntries int
	pool       sync.Pool

	mu    *sync.Mutex
	ll    *list.List
	cache map[LRU4Key]*list.Element
}

type entry struct {
	key   LRU4Key
	value interface{}
}

func New(maxEntries int) *LRU4Cache {
	return &LRU4Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[LRU4Key]*list.Element),
		pool: sync.Pool{
			New: func() interface{} {
				return &entry{}
			},
		},
	}
}

func (l *LRU4Cache) Get(key LRU4Key) (val interface{}, hit bool) {
	if l.cache == nil {
		return
	}
	if ele, ok := l.cache[key]; ok {
		l.promote(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

func (l *LRU4Cache) Add(key LRU4Key, val interface{}) {
	if l.cache == nil {
		l.cache = make(map[LRU4Key]*list.Element)
		l.ll = list.New()
	}
	if ee, ok := l.cache[key]; ok {
		l.promote(ee)
		ee.Value.(*entry).value = val
		return
	}

	et := l.pool.Get().(*entry)
	et.key, et.value = key, val
	ele := l.ll.PushFront(et)
	l.cache[key] = ele
	if l.MaxEntries != 0 && l.ll.Len() > l.MaxEntries {
		l.RemoveOldest()
	}
}

func (l *LRU4Cache) removeElement(e *list.Element) {
	l.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(l.cache, kv.key)
	l.pool.Put(e.Value.(*entry))
}

func (l *LRU4Cache) promote(e *list.Element) {
	if e == l.ll.Front() {
		return
	}
	p := e.Prev()
	for i := 0; i < defaultStep && p != l.ll.Front(); i++ {
		p = p.Prev()
	}
}

func (l *LRU4Cache) RemoveOldest() {
	if l.cache == nil {
		return
	}
	ele := l.ll.Back()
	if ele != nil {
		l.removeElement(ele)
	}
}

func (l *LRU4Cache) Remove(key LRU4Key) {
	if l.cache == nil {
		return
	}
	if ele, ok := l.cache[key]; ok {
		l.removeElement(ele)
	}
}
