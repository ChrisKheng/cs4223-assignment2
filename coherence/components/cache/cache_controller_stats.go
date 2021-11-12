package cache

type CacheControllerStats struct {
	NumAccessesToPrivateData int
	NumAccessesToSharedData  int
	NumCacheMisses           int
	NumCacheAccesses         int // Hit + miss
}
