/*
Package xact implements a Transaction struct to simulate a Bus transaction.
*/
package xact

import "time"

type Transaction struct {
	TransactionType   TransactionType
	Address           uint32
	RequestedDataSize uint32 // In words
	SendDataSize      uint32 // Only set this if you want to send a block from a cache to another cache
	SenderId          int    // MUST specify
}

type TransactionType int

const (
	Nil TransactionType = iota // Nil value for transaction since struct can't be nil in Go.
	BusRead
	BusReadX
	BusUpgr
	MemReadDone
	MemWriteDone
	FlushOpt
	Flush
	BusUpd
	UpdateDone
)

type ReleaseBus func()
type OnRequestGrantedCallBack func(timestamp time.Time) Transaction
type SnoopingCallBack func(transaction Transaction)
type HasCopyCallBack func(address uint32) bool
