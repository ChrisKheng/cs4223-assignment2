package bus

import "github.com/chriskheng/cs4223-assignment2/coherence/components/memory"

type Bus struct {
	memory                memory.Memory
	state                 BusState
	onRequestGrantedFuncs []func()
}

type BusState int

const (
	Ready BusState = iota
)

func (b *Bus) Execute() {
	switch b.state {
	case Ready:
		if len(b.onRequestGrantedFuncs) == 0 {
			return
		}
		b.onRequestGrantedFuncs[0]()
		b.onRequestGrantedFuncs = b.onRequestGrantedFuncs[1:]
	}
}

func (b *Bus) RequestAccess(onRequestGranted func()) {
	b.onRequestGrantedFuncs = append(b.onRequestGrantedFuncs, onRequestGranted)
}
