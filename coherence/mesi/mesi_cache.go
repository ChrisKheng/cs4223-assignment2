package mesi

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type MesiCache struct {
	*cache.BaseCache
}

func NewMesiCache(blockSize, associativity, cacheSize int) *MesiCache {
	return &MesiCache{
		BaseCache: cache.NewBaseCache(blockSize, associativity, cacheSize),
	}
}

func (c *MesiCache) Execute() {
}

func (c *MesiCache) RequestRead(address uint32) {

}

func (c *MesiCache) RequestWrite(address uint32) {

}
