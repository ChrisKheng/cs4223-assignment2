package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type BaseCacheController struct {
	bus         bus.Bus
	cacheClient CacheClient
	cacheDs     CacheDs
}

func NewBaseCache(blockSize, associativity, cacheSize int) *BaseCacheController {
	return &BaseCacheController{
		cacheDs: NewCacheDs(blockSize, associativity, cacheSize),
	}
}
