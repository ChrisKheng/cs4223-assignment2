package simulator

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
)

type BaseSimulator struct {
	Cores  []core.Core
	Bus    bus.Bus
	Memory memory.Memory
}
