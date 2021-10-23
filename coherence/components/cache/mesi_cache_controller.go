package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type MesiCacheController struct {
	*BaseCacheController
}

func NewMesiCache(bus *bus.Bus, blockSize, associativity, cacheSize int) *MesiCacheController {
	mesiCC := &MesiCacheController{
		BaseCacheController: NewBaseCache(bus, blockSize, associativity, cacheSize),
	}

	// TODO: Include this in NewDragonCache function too.
	bus.RegisterSnoopingCallBack(mesiCC.OnSnoop)
	bus.RegisterGatherReplyCallBack(mesiCC.ReceiveReplyCallBack)
	return mesiCC
}

func (cc *MesiCacheController) RequestRead(address uint32, callback func()) {
	cc.onClientRequestComplete = callback
	if cc.cache.Contain(address) {
		cc.state = CacheHit
	} else {
		cc.state = CacheMiss
		cc.currentTransaction = xact.Transaction{
			TransactionType:   xact.BusRead,
			Address:           address,
			Callback:          cc.OnReadComplete,
			RequestedDataSize: cc.cache.blockSizeInWords,
		}
	}
}

func (cc *MesiCacheController) RequestWrite(address uint32, callback func()) {

}

func (cc *MesiCacheController) OnReadComplete(reply xact.ReplyMsg) {
	if cc.state != WaitForPropagation {
		panic(fmt.Sprintf("onReadComplete of cache is called when cache is in %d state", cc.state))
	}

	cc.cache.Insert(cc.currentTransaction.Address)

	cc.onClientRequestComplete()
	cc.state = Ready

	cc.bus.ReleaseBus(cc.busAcquiredTimestamp)
}

func (cc *MesiCacheController) OnSnoop(transaction xact.Transaction) {
	if !cc.cache.Contain(transaction.Address) {
		return
	}

	// Should ignore FlushOpt
	switch transaction.TransactionType {
	case xact.BusRead:
		// TODO: Check if the cache state is not shared state. If it's shared state, don't have to do anything
		// Also need to check if it is in M state.
		// This assumes that the cache state is in Exclusive.
		cc.transactionToSendWhenReplying = xact.Transaction{
			TransactionType: xact.FlushOpt,
			Address:         transaction.Address,
			SendDataSize:    transaction.RequestedDataSize,
		}
		cc.needToReply = true
	}
}

func (cc *MesiCacheController) ReceiveReplyCallBack(replyCallback xact.ReplyCallback) {
	if !cc.needToReply {
		return
	}

	replyCallback(cc.transactionToSendWhenReplying, xact.ReplyMsg{})
	cc.needToReply = false
}
