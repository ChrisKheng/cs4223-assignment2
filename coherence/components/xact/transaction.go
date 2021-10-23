package xact

import "time"

type Transaction struct {
	TransactionType   TransactionType
	Address           uint32
	Callback          OnRequestCompleteCallBack
	RequestedDataSize uint32 // In words
	SendDataSize      uint32 // Only set this if you want to send a block from a cache to another cache
}

type TransactionType int

const (
	BusRead TransactionType = iota
	FlushOpt
	NoOp
)

type ReleaseBus func()
type OnRequestGrantedCallBack func(timestamp time.Time) Transaction
type OnRequestCompleteCallBack func(reply ReplyMsg)

type SnoopingCallBack func(transaction Transaction)
type GatherReplyCallBack func(replyCallback ReplyCallback)
type ReplyCallback func(transaction Transaction, replyMsg ReplyMsg)
