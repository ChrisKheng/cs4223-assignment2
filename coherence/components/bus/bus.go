package bus

import (
	"fmt"
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
	"github.com/chriskheng/cs4223-assignment2/coherence/constants"
)

const transferCycles int = 2 // Cycles needed to send a word from a cache to another. Must be at least 2.

type Bus struct {
	state                 BusState
	onRequestGrantedFuncs []xact.OnRequestGrantedCallBack
	snoopingCallBacks     []xact.SnoopingCallBack
	counter               int
	requestBeingProcessed xact.Transaction
	replyToSend           xact.Transaction
	busAcquiredTimestamp  time.Time
	stats                 BusStats
}

type BusState int

const (
	Ready BusState = iota
	ProcessingRequest
	RequestSent
	ProcessingReply
	ReplySent
)

type BusStats struct {
	DataTraffic      int
	NumInvalidations int
	NumUpdates       int
}

func NewBus() *Bus {
	bus := &Bus{
		state: Ready,
	}

	return bus
}

func (b *Bus) Execute() {
	switch b.state {
	case Ready:
		if len(b.onRequestGrantedFuncs) == 0 {
			return
		}
		b.busAcquiredTimestamp = time.Now()
		transaction := b.onRequestGrantedFuncs[0](b.busAcquiredTimestamp)
		b.requestBeingProcessed = transaction
		b.onRequestGrantedFuncs = b.onRequestGrantedFuncs[1:]

		b.transferDataAndRecordStats(transaction)
		b.state = ProcessingRequest
	case ProcessingRequest:
		b.counter--
		if b.counter <= 0 {
			for _, snoopingCallback := range b.snoopingCallBacks {
				snoopingCallback(b.requestBeingProcessed)
			}
			b.state = RequestSent
		}
	case ProcessingReply:
		b.counter--
		if b.counter <= 0 {
			b.state = ReplySent // MUST be before callback!
			for _, snoopingCallback := range b.snoopingCallBacks {
				snoopingCallback(b.replyToSend)
			}
		}
	}
}

func (b *Bus) ReleaseBus(timestamp time.Time) {
	if timestamp != b.busAcquiredTimestamp {
		panic("given timestamp to ReleaseBus() is not the same as busAcquiredTimestamp")
	}
	b.state = Ready
}

func (b *Bus) RegisterSnoopingCallBack(callback xact.SnoopingCallBack) {
	b.snoopingCallBacks = append(b.snoopingCallBacks, callback)
}

func (b *Bus) RequestAccess(onRequestGranted xact.OnRequestGrantedCallBack) {
	b.onRequestGrantedFuncs = append(b.onRequestGrantedFuncs, onRequestGranted)
}

func (b *Bus) Reply(transaction xact.Transaction) {
	if b.state == ProcessingReply {
		// Bus would only send the first reply it receives.
		return
	} else if !(b.state == RequestSent || b.state == ReplySent) {
		panic(fmt.Sprintf("bus's reply() is called when bus is in %d state\n", b.state))
	}

	b.transferDataAndRecordStats(transaction)
	b.recordStats(transaction)
	b.replyToSend = transaction
	b.state = ProcessingReply
}

func (b *Bus) GetStatistics() BusStats {
	return b.stats
}

func (b *Bus) transferDataAndRecordStats(transaction xact.Transaction) {
	// b.counter +1 to leave the send reply logic to Execute() cuz counter may be zero here if without +1.
	b.counter = transferCycles*(int(transaction.SendDataSize)) + 1
	b.recordStats(transaction)
}

func (b *Bus) recordStats(transaction xact.Transaction) {
	b.stats.DataTraffic += int(transaction.SendDataSize) * int(constants.WordSize)
	switch transaction.TransactionType {
	case xact.BusReadX, xact.BusUpgr:
		b.stats.NumInvalidations++
	}
}
