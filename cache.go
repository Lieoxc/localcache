package localcache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrShardCount    = errors.New("shard count must be power of two")
	ErrEntryNotFound = errors.New("entry not found")
	ErrkeyERR        = errors.New("key not equ")
)

type Cache struct {
	hashFunc    HashFunc
	bucketCount uint64
	bucketMask  uint64
	segments    []*segment
	locks       []sync.Mutex
	close       chan struct{}
}

func NewLocalCache(opts ...Opt) (*Cache, error) {
	options := &options{
		hashFunc:       NewDefaultHashFunc(),
		bucketCount:    defaultBucketCount,
		maxBytes:       defaultMaxBytes,
		cleanTime:      defaultCleanTIme,
		cleanupEnabled: true,
	}
	for _, each := range opts {
		each(options)
	}
	if !isPowerOfTwo(options.bucketCount) {
		return nil, ErrShardCount
	}
	segments := make([]*segment, options.bucketCount)
	maxSegmentBytes := (options.maxBytes + options.bucketCount - 1) / options.bucketCount
	for index := range segments {
		segments[index] = newSegment(maxSegmentBytes)
	}
	cache := &Cache{
		hashFunc:    options.hashFunc,
		bucketCount: options.bucketCount,
		bucketMask:  options.bucketCount - 1,
		segments:    segments,
		locks:       make([]sync.Mutex, options.bucketCount),
		close:       make(chan struct{}),
	}
	if options.cleanupEnabled {
		go cache.clean(options.cleanTime)
	}
	return cache, nil
}
func (c *Cache) clean(cleanTime time.Duration) {
	ticker := time.NewTicker(cleanTime)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			for index := 0; index < int(c.bucketCount); index++ {
				c.locks[index].Lock()
				c.segments[index].clean(t.Unix())
				c.locks[index].Unlock()
			}
		case <-c.close:
			return
		}
	}
}
func (c *Cache) Set(key string, value []byte) error {
	hashKey := c.hashFunc.Sum64(key)
	index := hashKey & c.bucketMask
	c.locks[index].Lock()
	defer c.locks[index].Unlock()
	c.segments[index].set(key, hashKey, value, defaultExpireTime)
	return nil
}
func (c *Cache) Del(key string) error {
	hashKey := c.hashFunc.Sum64(key)
	index := hashKey & c.bucketMask
	c.locks[index].Lock()
	defer c.locks[index].Unlock()
	err := c.segments[index].del(hashKey)
	return err
}
func (c *Cache) Get(key string) ([]byte, error) {
	hashKey := c.hashFunc.Sum64(key)
	index := hashKey & c.bucketMask
	c.locks[index].Lock()
	defer c.locks[index].Unlock()
	return c.segments[index].get(key, hashKey)
}
func (c *Cache) Len() int {
	cacheCount := 0
	for _, segment := range c.segments {
		cacheCount += segment.len()
	}
	return cacheCount
}
