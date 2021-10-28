package cache

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type DragonCacheController struct {
	*BaseCacheController
}

func NewDragonCache(id int, bus *bus.Bus, blockSize, associativity, cacheSize int) *DragonCacheController {
	dragonCC := &DragonCacheController{
		BaseCacheController: NewBaseCache(id, bus, blockSize, associativity, cacheSize),
	}

	bus.RegisterSnoopingCallBack(dragonCC.OnSnoop)
	bus.RegisterGatherReplyCallBack(dragonCC.ReceiveReplyCallBack)
	return dragonCC
}

func (cc *DragonCacheController) RequestRead(address uint32, callback func()) {
	cc.onClientRequestComplete = callback
	if cc.cache.Contain(address) {
		cc.state = ReadHit
	} else {
		cc.state = ReadMiss
	}
}

func (cc *DragonCacheController) RequestWrite(address uint32, callback func()) {

}

func (cc *DragonCacheController) OnReadComplete(reply xact.ReplyMsg) {

}

func (cc *DragonCacheController) OnReadExclusiveComplete(reply xact.ReplyMsg) {

}

func (cc *DragonCacheController) OnSnoop(transaction xact.Transaction) {

}

func (cc *DragonCacheController) ReceiveReplyCallBack(replyCallback xact.ReplyCallback) {

}
