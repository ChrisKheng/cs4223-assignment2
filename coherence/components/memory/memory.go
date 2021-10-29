package memory

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type Memory struct {
	bus                   *bus.Bus
	counter               int
	state                 MemoryState
	addressBeingProcessed uint32
	dataSizeInWords       uint32
	id                    int
}

type MemoryState int

const (
	Ready MemoryState = iota
	PrepareToReadResult
	PrepareToWriteResult
	ReadingResult
	WritingResult
)

const memLatency = 20

func NewMemory(id int, bus *bus.Bus) *Memory {
	memory := &Memory{id: id, bus: bus}
	memory.bus.RegisterSnoopingCallBack(memory.OnSnoop)
	return memory
}

func (m *Memory) Execute() {
	switch m.state {
	case PrepareToReadResult:
		m.state = ReadingResult
		m.counter = memLatency - 1
	case PrepareToWriteResult:
		m.state = WritingResult
		m.counter = memLatency - 1
	case ReadingResult:
		m.counter--
		if m.counter <= 0 {
			transaction := xact.Transaction{
				TransactionType: xact.MemReadDone,
				Address:         m.addressBeingProcessed,
				SendDataSize:    m.dataSizeInWords,
				SenderId:        m.id,
			}
			m.bus.Reply(transaction)
			m.state = Ready
		}
	case WritingResult:
		m.counter--
		if m.counter <= 0 {
			transaction := xact.Transaction{
				TransactionType: xact.MemWriteDone,
				Address:         m.addressBeingProcessed,
				SenderId:        m.id,
			}
			m.bus.Reply(transaction)
			m.state = Ready
		}
	}
}

func (m *Memory) OnSnoop(transaction xact.Transaction) {
	if transaction.SenderId == m.id {
		return
	}

	// if transaction.TransactionType != xact.FlushOpt && m.state != Ready {
	// 	panic(fmt.Sprintf("memory is in %d state when bus read is received", m.state))
	// }

	switch transaction.TransactionType {
	case xact.BusRead, xact.BusReadX:
		m.addressBeingProcessed = transaction.Address
		m.dataSizeInWords = transaction.RequestedDataSize
		m.state = PrepareToReadResult
	case xact.FlushOpt:
		m.dataSizeInWords = 0
		m.state = Ready
	case xact.Flush:
		m.addressBeingProcessed = transaction.Address
		m.dataSizeInWords = 0
		m.state = PrepareToWriteResult
	}
}
