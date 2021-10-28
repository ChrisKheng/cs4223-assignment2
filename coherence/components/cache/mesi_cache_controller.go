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
	bus.RegisterGatherReplyCallBack(mesiCC.ReceiveReplyCallBack)
	return mesiCC
}

func (cc *MesiCacheController) RequestRead(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		cc.state = ReadHit
	} else {
		cc.state = ReadMiss
		cc.stats.NumCacheMisses++
		cc.currentTransaction = xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			Callback:          cc.OnReadComplete,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}
	}
}

func (cc *MesiCacheController) RequestWrite(address uint32, callback func()) {
	cc.prepareForRequest(address, callback)

	if cc.cache.Contain(address) {
		cc.state = WriteHit
	} else {
		cc.state = WriteMiss
		cc.stats.NumCacheMisses++
		cc.currentTransaction = xact.Transaction{
			TransactionType:   xact.BusReadX,
			Address:           address,
			Callback:          cc.OnReadExclusiveComplete,
			RequestedDataSize: cc.cache.blockSizeInWords,
			SenderId:          cc.id,
		}
	}
}

func (cc *MesiCacheController) OnReadComplete(reply xact.ReplyMsg) {
	if cc.state != WaitForPropagation {
		panic(fmt.Sprintf("onReadComplete of cache is called when cache is in %d state", cc.state))
	}

	// TODO: Handle the case where a modified cache block got evicted!
	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	if reply.IsFromMem {
		cc.cacheStates[absoluteIndex] = Exclusive
	} else {
		cc.cacheStates[absoluteIndex] = Shared
		cc.stats.NumAccessesToSharedData++
	}

	cc.state = ReadHit
}

func (cc *MesiCacheController) OnReadExclusiveComplete(reply xact.ReplyMsg) {
	if cc.state != WaitForPropagation {
		panic(fmt.Sprintf("onWriteComplete of cache is called when cache is in %d state", cc.state))
	}

	_, _, absoluteIndex := cc.cache.Insert(cc.currentTransaction.Address)

	if !reply.IsFromMem {
		// TODO: This assumes that other cache is in S or E state
		// If the cache is in M state, need to increment NumAccessesToPrivateData instead
		cc.stats.NumAccessesToSharedData++
	}

	cc.cacheStates[absoluteIndex] = Modified
	cc.state = WriteHit
}

func (cc *MesiCacheController) OnSnoop(transaction xact.Transaction) {
	if transaction.SenderId == cc.id || !cc.cache.Contain(transaction.Address) {
		return
	}

	absoluteIndex := cc.cache.GetIndexInArray(transaction.Address)

	switch cc.cacheStates[absoluteIndex] {
	case Exclusive:
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
				cc.cacheStates[absoluteIndex] = Shared
			} else {
				cc.invalidateCache(transaction.Address, absoluteIndex)
			}
		}
	case Shared:
		switch transaction.TransactionType {
		case xact.BusReadX:
			cc.invalidateCache(transaction.Address, absoluteIndex)
		}
	}
}

func (cc *MesiCacheController) ReceiveReplyCallBack(replyCallback xact.ReplyCallback) {
	if !cc.needToReply {
		return
	}

	replyCallback(cc.transactionToSendWhenReplying, xact.ReplyMsg{IsFromMem: false})
	cc.needToReply = false
}

func (cc *MesiCacheController) invalidateCache(address uint32, absoluteIndex int) {
	cc.cacheStates[absoluteIndex] = Invalid
	cc.cache.Evict(address)
}
