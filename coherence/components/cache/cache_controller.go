package cache

type CacheController interface {
	Execute()
	RequestRead(address uint32, callback func())
	RequestWrite(address uint32, callback func())
}
