package mesi

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type MesiCacheController struct {
	*cache.BaseCacheController
}

func NewMesiCache(blockSize, associativity, cacheSize int) *MesiCacheController {
	return &MesiCacheController{
		BaseCacheController: cache.NewBaseCache(blockSize, associativity, cacheSize),
	}
}

func (c *MesiCacheController) Execute() {
}

func (c *MesiCacheController) RequestRead(address uint32) {

}

func (c *MesiCacheController) RequestWrite(address uint32) {

}
