package dragon

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

type DragonSimulator struct {
	*simulator.BaseSimulator
}

func NewDragonSimulator(numCores int, inputFilePrefix string, cacheSize int, associativity int, blockSize int) *DragonSimulator {
	cores := []*core.Core{}
	memory := memory.NewMemory(numCores)
	bus := bus.NewBus(memory)

	for i := 0; i < numCores; i++ {
		cache := cache.NewDragonCache(i, bus, blockSize, associativity, cacheSize)
		cores = append(cores, core.NewCore(i, inputFilePrefix, cache))
	}

	return &DragonSimulator{BaseSimulator: simulator.NewBaseSimulator(cores, bus, memory)}
}
