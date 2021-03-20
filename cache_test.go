package keyservice

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCache struct {
	data        sync.Map
	storeStat   map[string]int
	loadStat    map[string]int
	loadHitStat map[string]int
}

func newTestCache() *testCache {
	return &testCache{
		storeStat:   make(map[string]int),
		loadStat:    make(map[string]int),
		loadHitStat: make(map[string]int),
	}
}

func (c *testCache) Store(id string, content []byte) error {
	c.storeStat[id] += 1
	c.data.Store(id, content)
	return nil
}

func (c *testCache) Load(id string) (content []byte, err error) {
	c.loadStat[id] += 1
	if v, ok := c.data.Load(id); ok {
		c.loadHitStat[id] += 1
		return v.([]byte), nil
	} else {
		return nil, ErrNotFound
	}
}

func TestCache(t *testing.T) {
	c := NewCache()
	err := c.Store("test-key-1", []byte("hello"))
	assert.Nil(t, err)
	r, err := c.Load("test-key-1")
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), r)
}
