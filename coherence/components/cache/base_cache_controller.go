/*
Package cache implements:
* a Cache struct to simulate a processor cache
* a CacheController type struct to simulate a processor cache controller.

CacheController has a base type called BaseCacheController. It is embedded in specialised CacheController type
developed for simulating a particular cache coherence protocol, e.g. MesiCacheController.
*/
package cache

import (
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type BaseCacheController struct {
	bus                            *bus.Bus
	cache                          Cache
	state                          CacheControllerState
	onClientRequestComplete        func()
	requestedAddress               uint32
	currentTransaction             xact.Transaction
	needToReply                    bool
	transactionToSendWhenReplying  xact.Transaction
	busAcquiredTimestamp           time.Time
	isHoldingBus                   bool
	id                             int
	stats                          CacheControllerStats
	updateAccessStatsCallback      UpdateAccessStatsCallback
	iter                           int
	xactToIssueAfterEvictWriteBack xact.Transaction
}

type CacheControllerState int

const (
	Ready CacheControllerState = iota
	CacheHit
	RequestForBus
	WaitForBus
	WaitForRequestToComplete
	WaitForWriteBack
	WaitForEvictWriteBack
)

func NewBaseCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *BaseCacheController {
	baseCacheController := &BaseCacheController{
		bus:   bus,
		cache: NewCacheDs(blockSize, associativity, cacheSize),
		id:    id,
	}

	bus.RegisterHasCopy(baseCacheController.HasCopy)

	return baseCacheController
}

func (cc *BaseCacheController) RegisterUpdateAccessStatsCallback(callback UpdateAccessStatsCallback) {
	cc.updateAccessStatsCallback = callback
}

func (cc *BaseCacheController) Execute() {
	if cc.needToReply {
		cc.bus.Reply(cc.transactionToSendWhenReplying)
		cc.needToReply = false
		cc.transactionToSendWhenReplying = xact.Transaction{TransactionType: xact.Nil}
	}

	cc.iter++
	cc.cache.currentCycle++

	switch cc.state {
	case CacheHit:
		cc.stats.NumCacheAccesses++
		cc.cache.Access(cc.requestedAddress)
		cc.onClientRequestComplete()
		if cc.isHoldingBus {
			cc.bus.ReleaseBus(cc.busAcquiredTimestamp)
			cc.isHoldingBus = false
		}

		cc.updateAccessStatsCallback(cc.requestedAddress)
		cc.currentTransaction = xact.Transaction{TransactionType: xact.Nil}
		cc.state = Ready
	case RequestForBus:
		cc.bus.RequestAccess(cc.OnBusAccessGranted)
		cc.state = WaitForBus
	}
}

func (cc *BaseCacheController) OnBusAccessGranted(timestamp time.Time) xact.Transaction {
	cc.busAcquiredTimestamp = timestamp
	cc.isHoldingBus = true

	if cc.currentTransaction.TransactionType != xact.Flush {
		cc.state = WaitForRequestToComplete
	} else {
		cc.state = WaitForEvictWriteBack
	}

	return cc.currentTransaction
}

// MUST call in RequestRead and RequestWrite
func (cc *BaseCacheController) prepareForRequest(address uint32, callback func()) {
	cc.onClientRequestComplete = callback
	cc.requestedAddress = address
}

func (cc *BaseCacheController) GetStats() CacheControllerStats {
	return cc.stats
}

func (cc *BaseCacheController) HasCopy(address uint32) bool {
	return cc.cache.Contain(address)
}
