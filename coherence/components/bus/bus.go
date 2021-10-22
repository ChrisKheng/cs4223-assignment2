package bus

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/xact"
)

type Bus struct {
	memory                *memory.Memory
	state                 BusState
	onRequestGrantedFuncs []xact.OnRequestGrantedCallBack
	snoopingCallBacks     []xact.SnoopingCallBack
}

type BusState int

const (
	Ready BusState = iota
	Acquired
)

type ReleaseBus func()

func NewBus(memory *memory.Memory) *Bus {
	bus := &Bus{
		memory: memory,
		state:  Ready,
	}

	bus.RegisterSnoopingCallBack(memory.OnSnoop)

	return bus
}

func (b *Bus) Execute() {
	switch b.state {
	case Ready:
		if len(b.onRequestGrantedFuncs) == 0 {
			return
		}
		transaction := b.onRequestGrantedFuncs[0]()
		b.state = Acquired
		b.onRequestGrantedFuncs = b.onRequestGrantedFuncs[1:]

		for _, snoopingCallback := range b.snoopingCallBacks {
			snoopingCallback(transaction)
		}
	}
}

func (b *Bus) ReleaseBus() {
	b.state = Ready
}

func (b *Bus) RegisterSnoopingCallBack(callback xact.SnoopingCallBack) {
	b.snoopingCallBacks = append(b.snoopingCallBacks, callback)
}

func (b *Bus) Reply(transaction xact.Transaction) {

}

func (b *Bus) RequestAccess(onRequestGranted xact.OnRequestGrantedCallBack) {
	b.onRequestGrantedFuncs = append(b.onRequestGrantedFuncs, onRequestGranted)
}
