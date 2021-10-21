package cache

type Cache interface {
	Execute()
	RequestRead(address uint32)
	RequestWrite(address uint32)
}
