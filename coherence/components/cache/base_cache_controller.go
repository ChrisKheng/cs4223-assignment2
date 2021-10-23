package cache

import (
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type BaseCacheController struct {
	bus                           *bus.Bus
	cache                         Cache
	state                         CacheControllerState
	onClientRequestComplete       func()
	currentTransaction            xact.Transaction
	needToReply                   bool
	transactionToSendWhenReplying xact.Transaction
	busAcquiredTimestamp          time.Time
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
		bus:   bus,
		cache: NewCacheDs(blockSize, associativity, cacheSize),
	}
}

func (cc *BaseCacheController) Execute() {
	switch cc.state {
	case CacheHit:
		cc.onClientRequestComplete()
		cc.state = Ready
	case CacheMiss:
		cc.bus.RequestAccess(cc.OnBusAccessGranted)
		cc.state = WaitForBus
	}
}

func (cc *BaseCacheController) OnBusAccessGranted(timestamp time.Time) xact.Transaction {
	cc.busAcquiredTimestamp = timestamp
	cc.state = WaitForPropagation
	return cc.currentTransaction
}
