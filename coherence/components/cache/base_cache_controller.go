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
	requestedAddress              uint32
	currentTransaction            xact.Transaction
	needToReply                   bool
	transactionToSendWhenReplying xact.Transaction
	busAcquiredTimestamp          time.Time
	isHoldingBus                  bool
	id                            int
	stats                         CacheControllerStats
}

type CacheControllerState int

const (
	Ready CacheControllerState = iota
	ReadHit
	ReadMiss
	WriteHit
	WriteMiss
	WaitForBus
	WaitForPropagation
)

func NewBaseCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *BaseCacheController {
	// TODO: Register bus snooping callback here by calling bus.RegisterSnoopingCallBack
	return &BaseCacheController{
		bus:   bus,
		cache: NewCacheDs(blockSize, associativity, cacheSize),
		id:    id,
	}
}

func (cc *BaseCacheController) Execute() {
	switch cc.state {
	case ReadHit, WriteHit:
		cc.cache.Access(cc.requestedAddress)
		cc.onClientRequestComplete()
		if cc.isHoldingBus {
			cc.bus.ReleaseBus(cc.busAcquiredTimestamp)
			cc.isHoldingBus = false
		}

		cc.state = Ready
	case ReadMiss, WriteMiss:
		cc.bus.RequestAccess(cc.OnBusAccessGranted)
		cc.state = WaitForBus
	}
}

func (cc *BaseCacheController) OnBusAccessGranted(timestamp time.Time) xact.Transaction {
	cc.busAcquiredTimestamp = timestamp
	cc.state = WaitForPropagation
	cc.isHoldingBus = true
	return cc.currentTransaction
}

func (cc *BaseCacheController) GetStats() CacheControllerStats {
	return cc.stats
}

func (cc *BaseCacheController) prepareForRequest(address uint32, callback func()) {
	cc.onClientRequestComplete = callback
	cc.requestedAddress = address
	cc.stats.NumCacheAccesses++
}
