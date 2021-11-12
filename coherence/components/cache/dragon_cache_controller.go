package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type DragonCacheController struct {
	*BaseCacheController
	cacheStates                    []DragonCacheState
	requestType                    RequestTypes
	needToSendBusUpdAfterWriteBack bool
}

type DragonCacheState int

const (
	DragonModified DragonCacheState = iota // Put Invalid as first state just in case we forgot to initialise the state.
	DragonExclusive
	DragonSharedClean
	DragonSharedModified
)

type RequestTypes int

const (
	DragonRequestRead RequestTypes = iota
	DragonRequestWrite
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
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		cc.state = CacheHit
	} else {
		cc.state = RequestForBus
		cc.requestType = DragonRequestRead
		cc.stats.NumCacheMisses++

		busReadXact := xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}

		isToBeEvicted, evictedAddress, index := cc.cache.GetAddressToBeEvicted(address)
		if !isToBeEvicted || (cc.cacheStates[index] != DragonModified &&
			cc.cacheStates[index] != DragonSharedModified) {
			cc.currentTransaction = busReadXact
		} else { // to be evicted && DragonModified || DragonSharedModified
			cc.xactToIssueAfterEvictWriteBack = busReadXact
			cc.currentTransaction = xact.Transaction{
				TransactionType: xact.Flush,
				Address:         evictedAddress,
				SendDataSize:    cc.cache.blockSizeInWords,
				SenderId:        cc.id,
			}
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
				SendDataSize:    cc.cache.blockSizeInWords,
				SenderId:        cc.id,
			}
		case DragonModified:
			cc.state = CacheHit
		default:
			panic(fmt.Sprintf("cache state is in %d when cache data structure contains the address", state))
		}
	} else {
		cc.state = RequestForBus
		cc.requestType = DragonRequestWrite
		cc.stats.NumCacheMisses++
		busReadXact := xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}

		isToBeEvicted, evictedAddress, index := cc.cache.GetAddressToBeEvicted(address)
		if !isToBeEvicted || (cc.cacheStates[index] != DragonModified &&
			cc.cacheStates[index] != DragonSharedModified) {
			cc.currentTransaction = busReadXact
		} else { // to be evicted && DragonModified || DragonSharedModified
			cc.xactToIssueAfterEvictWriteBack = busReadXact
			cc.currentTransaction = xact.Transaction{
				TransactionType: xact.Flush,
				Address:         evictedAddress,
				SendDataSize:    cc.cache.blockSizeInWords,
				SenderId:        cc.id,
			}
		}
	}
}

func (cc *DragonCacheController) OnSnoop(transaction xact.Transaction) {
	switch cc.state {
	case WaitForEvictWriteBack:
		cc.handleSnoopWaitForEvictWriteBack(transaction)
	case WaitForRequestToComplete:
		cc.handleSnoopWaitForRequestToComplete(transaction)
	case WaitForWriteBack:
		cc.handleSnoopWriteBack(transaction)
	default:
		cc.handleSnoopOtherCases(transaction)
	}
}

func (cc *DragonCacheController) handleSnoopWaitForEvictWriteBack(transaction xact.Transaction) {
	if transaction.SenderId == cc.id {
		return
	}

	switch transaction.TransactionType {
	case xact.MemWriteDone:
		if transaction.Address != cc.currentTransaction.Address {
			panic("address evicted for write back is not the same as the address received for memwritedone")
		}

		cc.transactionToSendWhenReplying = cc.xactToIssueAfterEvictWriteBack
		cc.currentTransaction = cc.xactToIssueAfterEvictWriteBack
		cc.xactToIssueAfterEvictWriteBack = xact.Transaction{TransactionType: xact.Nil}
		cc.needToReply = true
		cc.state = WaitForRequestToComplete
	default:
		panic(fmt.Sprintf("Xact of type %d is received when cache controller is waiting for evict write back",
			transaction.TransactionType))
	}
}

func (cc *DragonCacheController) handleSnoopWaitForRequestToComplete(transaction xact.Transaction) {
	hasCopy := cc.bus.CheckHasCopy(cc.currentTransaction.Address)

	if transaction.SenderId == cc.id {
		if cc.currentTransaction.TransactionType == xact.BusUpd {
			cc.state = CacheHit
			absoluteIndex := cc.cache.GetIndexInArray(cc.currentTransaction.Address)
			if hasCopy {
				cc.cacheStates[absoluteIndex] = DragonSharedModified
			} else {
				cc.cacheStates[absoluteIndex] = DragonModified
			}
		}
		return
	}

	if !cc.cache.isSamePrefix(transaction.Address, cc.currentTransaction.Address) {
		panic("prefix of address received by cache controller is different than the prefix of the requested address while waiting for read to complete")
	}

	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	switch cc.currentTransaction.TransactionType {
	case xact.BusRead:
		if hasCopy {
			cc.stats.NumAccessesToSharedData++
			if cc.requestType == DragonRequestRead {
				cc.cacheStates[absoluteIndex] = DragonSharedClean
			} else {
				cc.cacheStates[absoluteIndex] = DragonSharedModified
				cc.needToSendBusUpdAfterWriteBack = true
			}
		} else {
			if cc.requestType == DragonRequestRead {
				cc.cacheStates[absoluteIndex] = DragonExclusive
			} else {
				cc.cacheStates[absoluteIndex] = DragonModified
			}
		}

		// Sender must be other cache
		switch transaction.TransactionType {
		case xact.Flush:
			cc.state = WaitForWriteBack
		case xact.MemReadDone:
			// TODO: May need to send busUpd here also
			cc.state = CacheHit
		default:
			// panic(fmt.Sprintf("transaction of type %s was received when cache controller is waiting for BusRead result", transaction.TransactionType))
		}
	}
}

func (cc *DragonCacheController) handleSnoopWriteBack(transaction xact.Transaction) {
	if transaction.TransactionType != xact.MemWriteDone {
		panic(fmt.Sprintf("transaction of type %d is received when cache controller %d is waiting for writeback, sender id: %d", transaction.TransactionType, cc.id, transaction.SenderId))
	} else if !cc.cache.isSamePrefix(transaction.Address, cc.currentTransaction.Address) {
		panic("tag of address written is not equal to the tag of address requested by cache controller")
	}

	if cc.needToSendBusUpdAfterWriteBack {
		cc.state = WaitForRequestToComplete
		cc.transactionToSendWhenReplying = xact.Transaction{
			TransactionType: xact.BusUpd,
			Address:         cc.currentTransaction.Address,
			SendDataSize:    cc.cache.blockSizeInWords,
			SenderId:        cc.id,
		}
		cc.needToReply = true
		cc.currentTransaction = cc.transactionToSendWhenReplying
		cc.needToSendBusUpdAfterWriteBack = false
	} else {
		cc.state = CacheHit
	}
}

func (cc *DragonCacheController) handleSnoopOtherCases(transaction xact.Transaction) {
	if transaction.SenderId == cc.id || !cc.cache.Contain(transaction.Address) {
		return
	}

	absoluteIndex := cc.cache.GetIndexInArray(transaction.Address)

	switch transaction.TransactionType {
	case xact.BusRead:
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
			panic("handleSnoopOtherCases BusRead undefined cacheStates")
		}
	case xact.BusUpd:
		switch cc.cacheStates[absoluteIndex] {
		case DragonSharedClean, DragonSharedModified:
			cc.cacheStates[absoluteIndex] = DragonSharedClean
			cc.stats.NumCacheUpdates++
		}
	}
}
