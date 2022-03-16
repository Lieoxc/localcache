package cache

import (
	"container/list"
	"errors"
	"sync"
)

type LRUCache struct {
	cap   int
	l     *list.List
	cache map[interface{}]*list.Element
	mux   *sync.RWMutex
}
type LruNode struct {
	value interface{}
	key   interface{}
}

func NewLRU(capacity int) *LRUCache {
	return &LRUCache{
		cap:   capacity,
		l:     new(list.List),
		cache: make(map[interface{}]*list.Element),
		mux:   new(sync.RWMutex),
	}
}

func (lru *LRUCache) Get(key interface{}) (val interface{}, ok bool) {
	if lru.l == nil {
		return
	}
	lru.mux.RLock()
	defer lru.mux.RUnlock()

	if v, ok := lru.cache[key]; ok {
		lru.l.MoveToFront(v)
		return v.Value.(*LruNode).value, true
	}
	return
}
func (lru *LRUCache) GetALL() []*LruNode {
	lru.mux.Lock()
	defer lru.mux.Unlock()
	var data []*LruNode
	for _, v := range lru.cache {
		data = append(data, v.Value.(*LruNode))
	}
	return data

}
func (lru *LRUCache) Add(key, val interface{}) error {
	if lru.l == nil {
		return errors.New("not init LRU list")
	}
	lru.mux.Lock()
	defer lru.mux.Unlock()
	//exist,replace it and move to front
	if v, ok := lru.cache[key]; ok {
		v.Value.(*LruNode).value = val
		lru.l.MoveToFront(v)
		return nil
	}
	//not exist,create it and move to front
	LruNodeElem := &LruNode{
		key:   key,
		value: val,
	}
	l := lru.l.PushFront(LruNodeElem)
	lru.cache[key] = l

	//check cap and delete Least  recently  use elem
	if lru.cap != 0 && lru.l.Len() > lru.cap {
		if e := lru.l.Back(); e != nil {
			lru.l.Remove(e)
			n := e.Value.(*LruNode)
			delete(lru.cache, n.key)
		}
	}
	return nil
}

func (lru *LRUCache) Del(key interface{}) {
	if lru.l == nil {
		return
	}
	lru.mux.Lock()
	defer lru.mux.Unlock()

	if elem, ok := lru.cache[key]; ok {
		lru.l.Remove(elem)
		delete(lru.cache, key)
	}
}
