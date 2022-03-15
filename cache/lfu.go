package cache

import (
	"container/list"
	"errors"
	"sync"
)

type LFUCache struct {
	cap   int
	l     *list.List
	cache map[interface{}]*list.Element
	mux   *sync.RWMutex
}
type LfuNode struct {
	value interface{}
	key   interface{}
}

func NewLFU(capacity int) *LFUCache {
	return &LFUCache{
		cap:   capacity,
		l:     new(list.List),
		cache: make(map[interface{}]*list.Element),
		mux:   new(sync.RWMutex),
	}
}

func (lfu *LFUCache) Get(key interface{}) (val interface{}, ok bool) {
	if lfu.l == nil {
		return
	}
	lfu.mux.RLock()
	defer lfu.mux.RUnlock()

	if v, ok := lfu.cache[key]; ok {
		lfu.l.MoveToFront(v)
		return v.Value.(*LfuNode).value, true
	}
	return
}
func (lfu *LFUCache) GetALL() []*LfuNode {
	lfu.mux.Lock()
	defer lfu.mux.Unlock()
	var data []*LfuNode
	for _, v := range lfu.cache {
		data = append(data, v.Value.(*LfuNode))
	}
	return data

}
func (lfu *LFUCache) Add(key, val interface{}) error {
	if lfu.l == nil {
		return errors.New("not init lfu list")
	}
	lfu.mux.Lock()
	defer lfu.mux.Unlock()
	//exist,replace it and move to front
	if v, ok := lfu.cache[key]; ok {
		v.Value.(*LfuNode).value = val
		lfu.l.MoveToFront(v)
		return nil
	}
	//not exist,create it and move to front
	NodeElem := &LfuNode{
		key:   key,
		value: val,
	}
	l := lfu.l.PushFront(NodeElem)
	lfu.cache[key] = l

	//check cap and delete last no use elem
	if lfu.cap != 0 && lfu.l.Len() > lfu.cap {
		if e := lfu.l.Back(); e != nil {
			lfu.l.Remove(e)
			n := e.Value.(*LfuNode)
			delete(lfu.cache, n.key)
		}
	}
	return nil
}

func (lfu *LFUCache) Del(key interface{}) {
	if lfu.l == nil {
		return
	}
	lfu.mux.Lock()
	defer lfu.mux.Unlock()

	if elem, ok := lfu.cache[key]; ok {
		lfu.l.Remove(elem)
		delete(lfu.cache, key)
	}
}
