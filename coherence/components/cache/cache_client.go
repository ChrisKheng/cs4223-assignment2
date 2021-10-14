package cache

type CacheClient interface {
	OnRequestComplete()
}
