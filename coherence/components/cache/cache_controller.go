package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/xact"

type CacheController interface {
	Execute()
	RequestRead(address uint32, callback func())
	RequestWrite(address uint32, callback func())
	OnSnoop(transaction xact.Transaction)
	HasCopy(address uint32) bool
	GetStats() CacheControllerStats
}
