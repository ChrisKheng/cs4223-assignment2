package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type DragonCacheController struct {
	*BaseCacheController
}

func NewDragonCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *DragonCacheController {
	return &DragonCacheController{
		BaseCacheController: NewBaseCache(id, bus, blockSize, associativity, cacheSize),
	}
}

func (c *DragonCacheController) RequestRead(address uint32, callback func()) {
	c.onClientRequestComplete = callback
	if c.cache.Contain(address) {
		c.state = ReadHit
	} else {
		c.state = ReadMiss
	}
}

func (c *DragonCacheController) RequestWrite(address uint32, callback func()) {

}
