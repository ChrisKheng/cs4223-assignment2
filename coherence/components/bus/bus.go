package bus

import (
	"fmt"
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

const transferCycles uint = 2 // Cycles needed to send a word from a cache to another. Must be at least 2.

type Bus struct {
	memory                *memory.Memory
	state                 BusState
	onRequestGrantedFuncs []xact.OnRequestGrantedCallBack
	snoopingCallBacks     []xact.SnoopingCallBack
	gatherReplyCallBacks  []xact.GatherReplyCallBack
	counter               int
	requestBeingProcessed xact.Transaction
	replyMsg              xact.ReplyMsg
	busAcquiredTimestamp  time.Time
	stats                 BusStats
}

type BusState int

const (
	Ready BusState = iota
	ProcessingRequest
	RequestSent
	WaitingForReply
	ProcessingReply
	ReplySent
)

type BusStats struct {
	DataTraffic int64
}

func NewBus(memory *memory.Memory) *Bus {
	bus := &Bus{
		memory: memory,
		state:  Ready,
	}

	bus.RegisterSnoopingCallBack(memory.OnSnoop)
	bus.RegisterGatherReplyCallBack(memory.ReceiveReplyCallBack)

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
		b.state = ProcessingRequest
	case ProcessingRequest:
		for _, snoopingCallback := range b.snoopingCallBacks {
			snoopingCallback(b.requestBeingProcessed)
		}

		b.state = RequestSent
	case RequestSent:
		b.state = WaitingForReply

		for _, gatherReplyCallback := range b.gatherReplyCallBacks {
			gatherReplyCallback(b.reply)
		}
	case ProcessingReply:
		b.counter--
		if b.counter <= 0 {
			b.state = ReplySent // MUST be before callback!
			b.requestBeingProcessed.Callback(b.replyMsg)
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

func (b *Bus) RegisterGatherReplyCallBack(callback xact.GatherReplyCallBack) {
	b.gatherReplyCallBacks = append(b.gatherReplyCallBacks, callback)
}

func (b *Bus) RequestAccess(onRequestGranted xact.OnRequestGrantedCallBack) {
	b.onRequestGrantedFuncs = append(b.onRequestGrantedFuncs, onRequestGranted)
}

func (b *Bus) reply(transaction xact.Transaction, reply xact.ReplyMsg) {
	// TODO: Remove b.state == ProcessingReply check after you have modified the reply
	// to BusRead logic such that only Exclusive state sends reply
	if !(b.state == WaitingForReply || b.state == ProcessingReply) {
		panic(fmt.Sprintf("bus's reply() is called when bus is in %d state\n", b.state))
	}

	// b.counter +1 to leave the send reply logic to Execute() cuz counter may be zero here if without +1.
	b.stats.DataTraffic += int64(transaction.SendDataSize)
	b.counter = int(2*transaction.SendDataSize) + 1

	// Snooping callback shouldn't sent to sender!
	// Probably can include sender index in transaction
	// Need to send reply to everyone else, e.g. to memory so that memory transit from
	// preparingToRead to ready in the case of cache-to-cache sharing.
	for _, snoopingCallback := range b.snoopingCallBacks {
		snoopingCallback(transaction)
	}

	b.replyMsg = reply
	b.state = ProcessingReply
}
