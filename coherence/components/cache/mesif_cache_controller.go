package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type MesifCacheController struct {
	*BaseCacheController
	xactToIssueAfterEvictWriteBack xact.Transaction
	cacheStates                    []mesifCacheState
}

type mesifCacheState int

const (
	mesifInvalid mesifCacheState = iota // Put Invalid as first state just in case we forgot to initialise the state.
	mesifModified
	mesifExclusive
	mesifShared
	mesifForward
)

func NewMesifCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *MesifCacheController {
	mesifCC := &MesifCacheController{
		BaseCacheController: NewBaseCache(id, bus, blockSize, associativity, cacheSize),
	}

	mesifCC.cacheStates = make([]mesifCacheState, len(mesifCC.cache.cacheArray))
	for i := range mesifCC.cacheStates {
		mesifCC.cacheStates[i] = mesifInvalid
	}

	bus.RegisterSnoopingCallBack(mesifCC.OnSnoop)
	return mesifCC
}

func (cc *MesifCacheController) RequestRead(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		cc.state = CacheHit
	} else {
		cc.state = RequestForBus
		cc.stats.NumCacheMisses++
		busReadXact := xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}

		isToBeEvicted, evictedAddress, index := cc.cache.GetAddressToBeEvicted(address)
		if !isToBeEvicted || cc.cacheStates[index] != mesifModified {
			cc.currentTransaction = busReadXact
		} else {
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

func (cc *MesifCacheController) RequestWrite(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		index := cc.cache.GetIndexInArray(address)
		state := cc.cacheStates[index]
		switch state {
		case mesifModified:
			cc.state = CacheHit
		case mesifExclusive:
			cc.state = CacheHit
			cc.cacheStates[index] = mesifModified
		case mesifShared:
			cc.state = RequestForBus
			cc.currentTransaction = xact.Transaction{
				TransactionType: xact.BusUpgr,
				Address:         address,
				SenderId:        cc.id,
			}
		default:
			panic(fmt.Sprintf("cache state is in %d when cache data structure contains the address", state))
		}
	} else {
		cc.state = RequestForBus
		cc.stats.NumCacheMisses++
		busReadXXact := xact.Transaction{
			TransactionType:   xact.BusReadX,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}

		isToBeEvicted, evictedAddress, index := cc.cache.GetAddressToBeEvicted(address)
		if !isToBeEvicted || cc.cacheStates[index] != mesifModified {
			cc.currentTransaction = busReadXXact
		} else {
			cc.xactToIssueAfterEvictWriteBack = busReadXXact
			cc.currentTransaction = xact.Transaction{
				TransactionType: xact.Flush,
				Address:         evictedAddress,
				SendDataSize:    cc.cache.blockSizeInWords,
				SenderId:        cc.id,
			}
		}
	}
}

func (cc *MesifCacheController) OnSnoop(transaction xact.Transaction) {
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

func (cc *MesifCacheController) handleSnoopWaitForEvictWriteBack(transaction xact.Transaction) {
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

func (cc *MesifCacheController) handleSnoopWaitForRequestToComplete(transaction xact.Transaction) {
	// Handle S -> M state
	if transaction.SenderId == cc.id && transaction.TransactionType == xact.BusUpgr {
		// Should have the same address (since the message is from the current sender itself (loopback))
		if transaction.Address != cc.currentTransaction.Address {
			panic("cache controller receives a different address than the address in BusUpgr")
		}
		cc.state = CacheHit
		index := cc.cache.GetIndexInArray(cc.currentTransaction.Address)
		if index == -1 {
			panic(fmt.Sprintf("index returned is -1, current iter %d", cc.iter))
		}

		cc.cacheStates[index] = mesifModified
		return
	}

	if transaction.SenderId == cc.id {
		return
	}

	// only for reply cases
	if !cc.cache.isSameTag(transaction.Address, cc.currentTransaction.Address) {
		fmt.Printf("iter: %d\n", cc.iter)
		panic("tag of address received by cache controller is different than the tag of the requested address while waiting for read to complete")
	}

	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	switch cc.currentTransaction.TransactionType {
	case xact.BusRead:
		if transaction.TransactionType == xact.Flush {
			cc.state = WaitForWriteBack
			cc.cacheStates[absoluteIndex] = mesifShared
			cc.stats.NumAccessesToPrivateData++
		} else if transaction.TransactionType == xact.MemReadDone || transaction.TransactionType == xact.FlushOpt {
			cc.state = CacheHit
			if transaction.TransactionType == xact.MemReadDone {
				cc.cacheStates[absoluteIndex] = mesifExclusive
			} else if transaction.TransactionType == xact.FlushOpt {
				cc.cacheStates[absoluteIndex] = mesifShared
				cc.stats.NumAccessesToSharedData++
			}
		} else {
			panic(fmt.Sprintf("transaction of type %d was received when cache controller is waiting for BusRead result", transaction.TransactionType))
		}
	case xact.BusReadX:
		cc.cacheStates[absoluteIndex] = mesifModified
		if transaction.TransactionType == xact.Flush {
			cc.state = WaitForWriteBack
			cc.stats.NumAccessesToPrivateData++
		} else if transaction.TransactionType == xact.MemReadDone || transaction.TransactionType == xact.FlushOpt {
			cc.state = CacheHit
			if transaction.TransactionType == xact.FlushOpt {
				cc.stats.NumAccessesToSharedData++
			}
		} else {
			panic(fmt.Sprintf("transaction of type %d was received when cache controller is waiting for BusReadX result", transaction.TransactionType))
		}
	}
}

// write to memory
func (cc *MesifCacheController) handleSnoopWriteBack(transaction xact.Transaction) {
	if transaction.TransactionType != xact.MemWriteDone {
		panic(fmt.Sprintf("transaction of type %d is received when cache controller %d is waiting for writeback, sender id: %d", transaction.TransactionType, cc.id, transaction.SenderId))
	} else if !cc.cache.isSameTag(transaction.Address, cc.currentTransaction.Address) {
		panic("tag of address written is not equal to the tag of address requested by cache controller")
	}

	cc.state = CacheHit
}

func (cc *MesifCacheController) handleSnoopOtherCases(transaction xact.Transaction) {
	if transaction.SenderId == cc.id || !cc.cache.Contain(transaction.Address) {
		return
	}

	absoluteIndex := cc.cache.GetIndexInArray(transaction.Address)

	switch cc.cacheStates[absoluteIndex] {
	case mesifModified:
		switch transaction.TransactionType {
		case xact.BusRead, xact.BusReadX:
			cc.transactionToSendWhenReplying = xact.Transaction{
				TransactionType: xact.Flush,
				Address:         transaction.Address,
				SendDataSize:    transaction.RequestedDataSize,
				SenderId:        cc.id,
			}
			cc.needToReply = true
			if transaction.TransactionType == xact.BusRead {
				cc.cacheStates[absoluteIndex] = mesifShared
			} else {
				cc.invalidateCache(transaction.Address, absoluteIndex)
			}

			// If the cache controller was waiting to flush and the address to flush is equal to
			// the address received in the snooped transaction, then the cache controller
			// don't have to flush when it got the ownership of the bus since it will flush
			// the cache line now.
			isWaitingToFlush := cc.state == WaitForBus &&
				cc.currentTransaction.TransactionType == xact.Flush &&
				cc.cache.isSameTag(cc.currentTransaction.Address, transaction.Address)
			if isWaitingToFlush {
				cc.currentTransaction = cc.xactToIssueAfterEvictWriteBack
				cc.xactToIssueAfterEvictWriteBack = xact.Transaction{TransactionType: xact.Nil}
			}
		}
	case mesifExclusive:
		switch transaction.TransactionType {
		case xact.BusRead, xact.BusReadX:
			cc.transactionToSendWhenReplying = xact.Transaction{
				TransactionType: xact.FlushOpt,
				Address:         transaction.Address,
				SendDataSize:    transaction.RequestedDataSize,
				SenderId:        cc.id,
			}
			cc.needToReply = true
			if transaction.TransactionType == xact.BusRead {
				cc.cacheStates[absoluteIndex] = mesifShared
			} else {
				cc.invalidateCache(transaction.Address, absoluteIndex)
			}
		}
	case mesifShared:
		switch transaction.TransactionType {
		case xact.BusReadX, xact.BusUpgr:
			needToChangeTransaction := cc.state == WaitForBus && cc.currentTransaction.TransactionType == xact.BusUpgr && cc.cache.isSameTag(cc.currentTransaction.Address, transaction.Address)
			if needToChangeTransaction {
				cc.currentTransaction = xact.Transaction{
					TransactionType:   xact.BusReadX,
					Address:           transaction.Address,
					RequestedDataSize: cc.cache.blockSizeInWords,
					SenderId:          cc.id,
				}
			}
			cc.invalidateCache(transaction.Address, absoluteIndex)
		}
	}
}

func (cc *MesifCacheController) invalidateCache(address uint32, absoluteIndex int) {
	cc.cacheStates[absoluteIndex] = mesifInvalid
	cc.cache.Evict(address)
}
