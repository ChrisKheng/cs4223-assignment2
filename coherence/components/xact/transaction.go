package xact

type Transaction struct {
	TransactionType TransactionType
	Address         uint32
	Callback        OnRequestCompleteCallBack
}

type TransactionType int

const (
	BusRead TransactionType = iota
	FlushOpt
)

type ReleaseBus func()
type OnRequestGrantedCallBack func() Transaction
type SnoopingCallBack func(transaction Transaction) (bool, Transaction)
type OnRequestCompleteCallBack func(reply ReplyMsg)
