package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type BaseCache struct {
	bus         bus.Bus
	cacheClient CacheClient
	cacheDs     CacheDs
}

func NewBaseCache(blockSize, associativity, cacheSize int) *BaseCache {
	return &BaseCache{
		cacheDs: NewCacheDs(blockSize, associativity, cacheSize),
	}
}
