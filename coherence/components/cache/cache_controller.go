package cache

import "github.com/chriskheng/cs4223-assignment2/coherence/components/xact"

type CacheController interface {
	Execute()
	RequestRead(address uint32, callback func())
	RequestWrite(address uint32, callback func())
	OnReadComplete(reply xact.ReplyMsg)
	OnWriteComplete(reply xact.ReplyMsg)
	OnSnoop(transaction xact.Transaction)
	ReceiveReplyCallBack(replyCallback xact.ReplyCallback)
}
