package memory

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type Memory struct {
	counter               int
	state                 MemoryState
	addressBeingProcessed uint32
	dataSizeInWords       uint32
	replyCallback         xact.ReplyCallback
}

type MemoryState int

const (
	Ready MemoryState = iota
	PrepareToReadResult
	ReadingResult
	WritingResult
)

const memLatency = 100

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Execute() {
	switch m.state {
	case PrepareToReadResult:
		m.state = ReadingResult
		m.counter = memLatency - 1
	case ReadingResult:
		m.counter--
		if m.counter <= 0 {
			m.replyCallback(xact.Transaction{TransactionType: xact.NoOp, Address: m.addressBeingProcessed, SendDataSize: m.dataSizeInWords}, xact.ReplyMsg{IsFromMem: true})
			m.state = Ready
			m.replyCallback = nil
		}
	}
}

func (m *Memory) OnSnoop(transaction xact.Transaction) {
	// if transaction.TransactionType != xact.FlushOpt && m.state != Ready {
	// 	panic(fmt.Sprintf("memory is in %d state when bus read is received", m.state))
	// }

	switch transaction.TransactionType {
	case xact.BusRead:
		m.addressBeingProcessed = transaction.Address
		m.dataSizeInWords = transaction.RequestedDataSize
		m.state = PrepareToReadResult
	case xact.FlushOpt:
		m.dataSizeInWords = 0
		m.replyCallback = nil
		m.state = Ready
	}
}

func (m *Memory) ReceiveReplyCallBack(replyCallback xact.ReplyCallback) {
	m.replyCallback = replyCallback
}
