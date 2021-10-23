package mesi

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/constants"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

type MesiSimulator struct {
	*simulator.BaseSimulator
}

func NewMesiSimulator(inputFilePrefix string, cacheSize int, associativity int, blockSize int) *MesiSimulator {
	cores := []*core.Core{}
	memory := memory.NewMemory(constants.NumCores)
	bus := bus.NewBus(memory)

	for i := 0; i < constants.NumCores; i++ {
		cache := cache.NewMesiCache(i, bus, blockSize, associativity, cacheSize)
		cores = append(cores, core.NewCore(i, inputFilePrefix, cache))
	}

	return &MesiSimulator{BaseSimulator: simulator.NewBaseSimulator(cores, bus, memory)}
}
