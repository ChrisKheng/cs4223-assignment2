package cache

type Cache interface {
	Execute()
	RequestRead()
	RequestWrite()
}
