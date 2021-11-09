package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type DragonCacheController struct {
	*BaseCacheController
	cacheStates []DragonCacheState
}

type DragonCacheState int

const (
	DragonModified DragonCacheState = iota // Put Invalid as first state just in case we forgot to initialise the state.
	DragonExclusive
	DragonSharedClean
	DragonSharedModified
)

func NewDragonCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *DragonCacheController {
	dragonCC := &DragonCacheController{
		BaseCacheController: NewBaseCache(id, bus, blockSize, associativity, cacheSize),
	}

	dragonCC.cacheStates = make([]DragonCacheState, len(dragonCC.cache.cacheArray))
	for i := range dragonCC.cacheStates {
		dragonCC.cacheStates[i] = DragonModified
	}

	bus.RegisterSnoopingCallBack(dragonCC.OnSnoop)
	return dragonCC
}

func (cc *DragonCacheController) RequestRead(address uint32, callback func()) {
	cc.onClientRequestComplete = callback
	cc.stats.NumCacheAccesses++
	if cc.cache.Contain(address) {
		cc.state = CacheHit
	} else {
		cc.state = RequestForBus
		cc.stats.NumCacheMisses++

		cc.currentTransaction = xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}
	}
}

func (cc *DragonCacheController) RequestWrite(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		index := cc.cache.GetIndexInArray(address)
		state := cc.cacheStates[index]

		switch state {
		case DragonExclusive:
			cc.state = CacheHit
			cc.cacheStates[index] = DragonModified
		case DragonSharedClean, DragonSharedModified:
			cc.state = RequestForBus
			cc.currentTransaction = xact.Transaction{
				TransactionType: xact.BusUpd,
				Address:         address,
				SenderId:        cc.id,
			}
		case DragonModified:
			cc.state = CacheHit
		default:
			panic(fmt.Sprintf("cache state is in %d when cache data structure contains the address", state))
		}
	} else {
		cc.state = RequestForBus
		cc.stats.NumCacheMisses++
		cc.currentTransaction = xact.Transaction{
			TransactionType: xact.BusRead,
			Address:         address,
			SenderId:        cc.id,
		}
	}
}

func (cc *DragonCacheController) OnSnoop(transaction xact.Transaction) {
	switch cc.state {
	case WaitForRequestToComplete:
		cc.handleSnoopWaitForRequestToComplete(transaction)
	case WaitForWriteBack:
		cc.handleSnoopWriteBack(transaction)
	default:
		cc.handleSnoopOtherCases(transaction)
	}
}

func (cc *DragonCacheController) handleSnoopWaitForRequestToComplete(transaction xact.Transaction) {
	if transaction.SenderId == cc.id {

		if cc.currentTransaction.TransactionType == xact.BusUpd {
			cc.state = CacheHit
		}
		return
	}

	if !cc.cache.isSamePrefix(transaction.Address, cc.currentTransaction.Address) {
		panic("prefix of address received by cache controller is different than the prefix of the requested address while waiting for read to complete")
	}

	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	switch cc.currentTransaction.TransactionType {
	case xact.BusRead:
		// Sender must be other cache
		switch transaction.TransactionType {
		case xact.Flush:
			cc.state = WaitForWriteBack
			cc.cacheStates[absoluteIndex] = DragonSharedClean
			cc.stats.NumAccessesToSharedData++
		case xact.MemReadDone:
			cc.state = CacheHit
			if transaction.TransactionType == xact.MemReadDone {
				cc.cacheStates[absoluteIndex] = DragonExclusive
			}

		default:
			// panic(fmt.Sprintf("transaction of type %s was received when cache controller is waiting for BusRead result", transaction.TransactionType))
		}
	}
}

func (cc *DragonCacheController) handleSnoopWriteBack(transaction xact.Transaction) {
	cc.state = CacheHit
}

func (cc *DragonCacheController) handleSnoopOtherCases(transaction xact.Transaction) {

	if transaction.SenderId == cc.id {
		return
	}

	absoluteIndex := cc.cache.GetIndexInArray(transaction.Address)

	switch transaction.TransactionType {
	case xact.BusRead:
		if !cc.cache.Contain(transaction.Address) {
			return
		}

		switch cc.cacheStates[absoluteIndex] {
		case DragonExclusive, DragonSharedClean:
			cc.needToReply = false
			cc.cacheStates[absoluteIndex] = DragonSharedClean

		case DragonSharedModified, DragonModified:
			cc.transactionToSendWhenReplying = xact.Transaction{
				TransactionType: xact.Flush,
				Address:         transaction.Address,
				SendDataSize:    transaction.RequestedDataSize,
				SenderId:        cc.id,
			}
			cc.needToReply = true
			cc.cacheStates[absoluteIndex] = DragonSharedModified

		default:
			panic("handleSnoopOtherCases BusRead undefine cacheStates")
		}
	case xact.BusUpd:
		if !cc.cache.Contain(transaction.Address) {
			return
		}

		switch cc.cacheStates[absoluteIndex] {
		case DragonSharedClean, DragonSharedModified:
			cc.cacheStates[absoluteIndex] = DragonSharedClean
			cc.stats.NumCacheUpdates++
		}
	}

}
