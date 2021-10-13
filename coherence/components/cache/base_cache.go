package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/bus"

type BaseCache struct {
	bus         bus.Bus
	cacheClient CacheClient
}
