package cache

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type BaseCacheController struct {
	bus                     *bus.Bus
	cacheDs                 Cache
	state                   CacheControllerState
	onClientRequestComplete func()
	currentTransaction      xact.Transaction
}

type CacheControllerState int

const (
	Ready CacheControllerState = iota
	CacheHit
	CacheMiss
	WaitForBus
	WaitForPropagation
)

func NewBaseCache(bus *bus.Bus, blockSize, associativity, cacheSize int) *BaseCacheController {
	// TODO: Register bus snooping callback here by calling bus.RegisterSnoopingCallBack
	return &BaseCacheController{
		bus:     bus,
		cacheDs: NewCacheDs(blockSize, associativity, cacheSize),
	}
}

func (c *BaseCacheController) Execute() {
	switch c.state {
	case CacheHit:
		c.onClientRequestComplete()
		c.state = Ready
	case CacheMiss:
		c.bus.RequestAccess(c.OnBusAccessGranted)
		c.state = WaitForBus
	}
}

func (c *BaseCacheController) OnBusAccessGranted() xact.Transaction {
	c.state = WaitForPropagation
	return c.currentTransaction
}
