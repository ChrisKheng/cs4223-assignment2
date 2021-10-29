package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type MesiCacheController struct {
	*BaseCacheController
	cacheStates []MesiCacheState
}

type MesiCacheState int

const (
	Invalid MesiCacheState = iota // Put Invalid as first state just in case we forgot to initialise the state.
	Modified
	Exclusive
	Shared
)

func NewMesiCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *MesiCacheController {
	mesiCC := &MesiCacheController{
		BaseCacheController: NewBaseCache(id, bus, blockSize, associativity, cacheSize),
	}

	mesiCC.cacheStates = make([]MesiCacheState, len(mesiCC.cache.cacheArray))
	for i := range mesiCC.cacheStates {
		mesiCC.cacheStates[i] = Invalid
	}

	bus.RegisterSnoopingCallBack(mesiCC.OnSnoop)
	return mesiCC
}

func (cc *MesiCacheController) RequestRead(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

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

func (cc *MesiCacheController) RequestWrite(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		index := cc.cache.GetIndexInArray(address)
		state := cc.cacheStates[index]
		switch state {
		case Modified:
			cc.state = CacheHit
		case Exclusive:
			cc.state = CacheHit
			cc.cacheStates[index] = Modified
		case Shared:
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
		cc.currentTransaction = xact.Transaction{
			TransactionType:   xact.BusReadX,
			Address:           address,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}
	}
}

func (cc *MesiCacheController) OnSnoop(transaction xact.Transaction) {
	switch cc.state {
	case WaitForRequestToComplete:
		cc.handleSnoopWaitForRequestToComplete(transaction)
	case WaitForWriteBack:
		cc.handleSnoopWriteBack(transaction)
	default:
		cc.handleSnoopOtherCases(transaction)
	}
}

func (cc *MesiCacheController) handleSnoopWaitForRequestToComplete(transaction xact.Transaction) {
	// TODO: Handle the case where a modified cache block got evicted!
	// Handle S -> M state
	if transaction.SenderId == cc.id && transaction.TransactionType == xact.BusUpgr {
		if transaction.Address != cc.currentTransaction.Address {
			panic("cache controller receives a different address than the address in BusUpgr")
		}
		cc.state = CacheHit
		index := cc.cache.GetIndexInArray(cc.currentTransaction.Address)
		if index == -1 {
			panic("index returned is -1")
		}

		cc.cacheStates[index] = Modified
		return
	}

	if transaction.SenderId == cc.id {
		return
	}

	if transaction.Address != cc.currentTransaction.Address {
		panic("cache controller receives a different address than the requested address while waiting for read to complete")
	}

	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	switch cc.currentTransaction.TransactionType {
	case xact.BusRead:
		if transaction.TransactionType == xact.Flush {
			cc.state = WaitForWriteBack
			cc.cacheStates[absoluteIndex] = Shared
			cc.stats.NumAccessesToPrivateData++
		} else if transaction.TransactionType == xact.MemReadDone || transaction.TransactionType == xact.FlushOpt {
			cc.state = CacheHit
			if transaction.TransactionType == xact.MemReadDone {
				cc.cacheStates[absoluteIndex] = Exclusive
			} else if transaction.TransactionType == xact.FlushOpt {
				cc.cacheStates[absoluteIndex] = Shared
				cc.stats.NumAccessesToSharedData++
			}
		} else {
			panic(fmt.Sprintf("transaction of type %d was received when cache controller is waiting for BusRead result", transaction.TransactionType))
		}
	case xact.BusReadX:
		cc.cacheStates[absoluteIndex] = Modified
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

func (cc *MesiCacheController) handleSnoopWriteBack(transaction xact.Transaction) {
	if transaction.TransactionType != xact.MemWriteDone {
		panic(fmt.Sprintf("transaction of type %d is received when cache controller %d is waiting for writeback, sender id: %d", transaction.TransactionType, cc.id, transaction.SenderId))
	} else if transaction.Address != cc.currentTransaction.Address {
		panic("address written is not equal to the address requested by cache controller")
	}

	cc.state = CacheHit
}

func (cc *MesiCacheController) handleSnoopOtherCases(transaction xact.Transaction) {
	if transaction.SenderId == cc.id || !cc.cache.Contain(transaction.Address) {
		return
	}

	absoluteIndex := cc.cache.GetIndexInArray(transaction.Address)

	switch cc.cacheStates[absoluteIndex] {
	case Modified:
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
				cc.cacheStates[absoluteIndex] = Shared
			} else {
				cc.invalidateCache(transaction.Address, absoluteIndex)
			}
		}
	case Exclusive:
		switch transaction.TransactionType {
		case xact.BusRead, xact.BusReadX:
			cc.transactionToSendWhenReplying = xact.Transaction{
				TransactionType: xact.FlushOpt,
				Address:         transaction.Address,
				SendDataSize:    transaction.RequestedDataSize,
				SenderId:        cc.id,
			}
			cc.needToReply = true // TODO: add this in execute()
			if transaction.TransactionType == xact.BusRead {
				cc.cacheStates[absoluteIndex] = Shared
			} else {
				cc.invalidateCache(transaction.Address, absoluteIndex)
			}
		}
	case Shared:
		switch transaction.TransactionType {
		case xact.BusReadX, xact.BusUpgr:
			needToChangeTransaction := cc.state == WaitForBus && cc.currentTransaction.TransactionType == xact.BusUpgr &&
				cc.currentTransaction.Address == transaction.Address
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

func (cc *MesiCacheController) invalidateCache(address uint32, absoluteIndex int) {
	cc.cacheStates[absoluteIndex] = Invalid
	cc.cache.Evict(address)
}
