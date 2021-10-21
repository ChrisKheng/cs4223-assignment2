package dragon

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type DragonCache struct {
	*cache.BaseCache
}

func NewDragonCache(blockSize, associativity, cacheSize int) *DragonCache {
	return &DragonCache{
		BaseCache: cache.NewBaseCache(blockSize, associativity, cacheSize),
	}
}

func (c *DragonCache) Execute() {
}

func (c *DragonCache) RequestRead(cache uint32) {

}

func (c *DragonCache) RequestWrite(cache uint32) {

}
