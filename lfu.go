package cache

import (
	"container/list"
)

type LFUCache struct {
	cap        int
	mincount   int
	keyToNode  map[interface{}]*LfuNode //save key :value and count
	countToKey map[interface{}]*list.List
	//mux        *sync.RWMutex
}
type LfuNode struct {
	value interface{}
	key   interface{}
	count int
}

func NewLFU(capacity int) *LFUCache {
	return &LFUCache{
		cap:        capacity,
		mincount:   0,
		keyToNode:  make(map[interface{}]*LfuNode),
		countToKey: make(map[interface{}]*list.List),
		//mux:        new(sync.RWMutex),
	}
}

func (lfu *LFUCache) Get(key interface{}) (val interface{}, ok bool) {
	// lfu.mux.RLock()
	// defer lfu.mux.RUnlock()
	if lfu.cap == 0 {
		return
	}

	if _, ok := lfu.keyToNode[key]; !ok {
		return -1, false
	}
	node := lfu.keyToNode[key]
	count := node.count

	for e := lfu.countToKey[count].Front(); e != nil; e = e.Next() {
		if e.Value.(*LfuNode).key == key {
			lfu.countToKey[key].Remove(e)
			break
		}
	}

	//当前count已经是最小值，且唯一
	if count == lfu.mincount && lfu.countToKey[count].Len() == 0 {
		//最小标记需要+1
		lfu.mincount++
	}

	//add to list
	node.count += 1
	if _, ok := lfu.countToKey[node.count]; !ok {
		lfu.countToKey[node.count] = new(list.List)
	}
	lfu.countToKey[node.count].PushBack(node)

	return node.value, true
}

func (lfu *LFUCache) Add(key, val interface{}) error {
	// lfu.mux.Lock()
	// defer lfu.mux.Unlock()

	if lfu.cap == 0 {
		return nil
	}
	if node, ok := lfu.keyToNode[key]; ok {
		node.value = val
		lfu.keyToNode[key] = node
		lfu.Get(key) //flush count
		return nil
	}
	// need delete other elem
	if len(lfu.keyToNode) == lfu.cap {
		removeNode := lfu.countToKey[lfu.mincount].Front()
		rmKey := removeNode.Value.(*LfuNode).key
		lfu.countToKey[lfu.mincount].Remove(removeNode)
		delete(lfu.keyToNode, rmKey)
	}
	// a new elem
	node := &LfuNode{
		key:   key,
		value: val,
		count: 1,
	}
	lfu.keyToNode[key] = node
	lfu.mincount = 1
	//add to list
	if _, ok := lfu.countToKey[node.count]; !ok {
		lfu.countToKey[node.count] = new(list.List)
	}
	lfu.countToKey[node.count].PushBack(node)
	return nil
}
