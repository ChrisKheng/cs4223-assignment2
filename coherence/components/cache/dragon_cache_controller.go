package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type DragonCacheController struct {
	*BaseCacheController
}

func NewDragonCache(bus *bus.Bus, blockSize, associativity, cacheSize int) *DragonCacheController {
	return &DragonCacheController{
		BaseCacheController: NewBaseCache(bus, blockSize, associativity, cacheSize),
	}
}

func (c *DragonCacheController) RequestRead(address uint32, callback func()) {
	c.onClientRequestComplete = callback
	if c.cache.Contain(address) {
		c.state = CacheHit
	} else {
		c.state = CacheMiss
	}
}

func (c *DragonCacheController) RequestWrite(address uint32, callback func()) {

}
