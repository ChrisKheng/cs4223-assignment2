package cache

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type MesiCacheController struct {
	*BaseCacheController
}

// Probably need to take in a list of hasCopy callback methods
func NewMesiCache(bus *bus.Bus, blockSize, associativity, cacheSize int) *MesiCacheController {
	return &MesiCacheController{
		BaseCacheController: NewBaseCache(bus, blockSize, associativity, cacheSize),
	}
}

func (c *MesiCacheController) RequestRead(address uint32, callback func()) {
	c.onClientRequestComplete = callback
	if c.cacheDs.Contain(address) {
		c.state = CacheHit
	} else {
		c.state = CacheMiss
		c.currentTransaction = xact.Transaction{
			TransactionType: xact.BusRead,
			Address:         address,
			Callback:        c.OnReadComplete,
		}
	}
}

func (c *MesiCacheController) RequestWrite(address uint32, callback func()) {

}

func (c *MesiCacheController) OnReadComplete(reply xact.ReplyMsg) {
	if c.state != WaitForPropagation {
		panic(fmt.Sprintf("onReadComplete of cache is called when cache is in %d state", c.state))
	}

	c.cacheDs.Insert(c.currentTransaction.Address)

	c.onClientRequestComplete()
	c.state = Ready

	// Temporary hack
	c.bus.ReleaseBus()
}
