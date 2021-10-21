package cache

type Cache interface {
	Execute()
	RequestRead(address uint32, callback func())
	RequestWrite(address uint32, callback func())
}
