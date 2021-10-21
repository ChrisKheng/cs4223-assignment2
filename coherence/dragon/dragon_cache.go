package dragon

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type DragonCacheController struct {
	*cache.BaseCacheController
}

func NewDragonCache(blockSize, associativity, cacheSize int) *DragonCacheController {
	return &DragonCacheController{
		BaseCacheController: cache.NewBaseCache(blockSize, associativity, cacheSize),
	}
}
