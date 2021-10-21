package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type BaseCacheController struct {
	bus               bus.Bus
	cacheDs           CacheDs
	state             CacheControllerState
	onRequestComplete func()
}

type CacheControllerState int

const (
	Ready CacheControllerState = iota
	CacheHit
	CacheMiss
)

func NewBaseCache(blockSize, associativity, cacheSize int) *BaseCacheController {
	return &BaseCacheController{
		cacheDs: NewCacheDs(blockSize, associativity, cacheSize),
	}
}

func (c *BaseCacheController) Execute() {
	if c.state == CacheHit {
		c.onRequestComplete()
		c.state = Ready
	} else if c.state == CacheMiss {
		// TODO: Get access to bus
		c.onRequestComplete()
		c.state = Ready
	}
	c.onRequestComplete = nil
}

func (c *BaseCacheController) RequestRead(address uint32, callback func()) {
	c.onRequestComplete = callback
	if c.cacheDs.Contain(address) {
		c.state = CacheHit
	} else {
		c.state = CacheMiss
	}
}

func (c *BaseCacheController) RequestWrite(address uint32, callback func()) {

}
