package memory

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type Memory struct {
	counter                   int
	state                     MemoryState
	onRequestCompleteCallback xact.OnRequestCompleteCallBack
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
		if m.counter == 0 {
			m.onRequestCompleteCallback(xact.ReplyMsg{IsFromMem: true})
			m.state = Ready
			m.onRequestCompleteCallback = nil
		}
	}
}

func (m *Memory) OnSnoop(transaction xact.Transaction) (bool, xact.Transaction) {
	// if transaction.TransactionType != xact.FlushOpt && m.state != Ready {
	// 	panic(fmt.Sprintf("memory is in %d state when bus read is received", m.state))
	// }

	switch transaction.TransactionType {
	case xact.BusRead:
		m.state = PrepareToReadResult
		m.onRequestCompleteCallback = transaction.Callback
	case xact.FlushOpt:
		m.state = Ready
		m.onRequestCompleteCallback = nil
	}

	return false, xact.Transaction{}
}
