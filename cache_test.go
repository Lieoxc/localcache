package localcache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type cacheTestSuite struct {
	suite.Suite
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(cacheTestSuite))
}

func (h *cacheTestSuite) SetupSuite() {}

func (h *cacheTestSuite) TestSetAndGet() {
	cache, err := NewLocalCache()
	assert.Equal(h.T(), nil, err)
	key := "lxc"
	value := []byte("mycache")

	err = cache.Set(key, value)
	assert.Equal(h.T(), nil, err)

	res, err := cache.Get(key)
	assert.Equal(h.T(), nil, err)
	assert.Equal(h.T(), value, res)
	h.T().Logf("get value is %s", string(res))
}

func (h *cacheTestSuite) TestLen() {
	cache, err := NewLocalCache()
	assert.Equal(h.T(), nil, err)

	value := []byte("公众号：Golang梦工厂")
	for index := 0; index < 1000; index++ {
		key := fmt.Sprintf("asong%03d", index)
		err = cache.Set(key, value)
		assert.Equal(h.T(), nil, err)
	}

	length := cache.Len()
	assert.Equal(h.T(), 1000, length)
	h.T().Logf("length == %d", length)
}

func (h *cacheTestSuite) TestDel() {
	cache, err := NewLocalCache()
	assert.Equal(h.T(), nil, err)
	key := "lxc"
	value := []byte("mycache")

	err = cache.Set(key, value)
	assert.Equal(h.T(), nil, err)

	res, err := cache.Get(key)
	assert.Equal(h.T(), nil, err)
	assert.Equal(h.T(), value, res)
	h.T().Logf("res == %s", string(res))
	err = cache.Del(key)
	assert.Equal(h.T(), nil, err)

	_, err = cache.Get(key)
	assert.Equal(h.T(), ErrEntryNotFound, err)
}
