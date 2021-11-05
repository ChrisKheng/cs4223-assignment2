package mesif

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/constants"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

type MesifSimulator struct {
	*simulator.BaseSimulator
}

func NewMesifSimulator(inputFilePrefix string, cacheSize int, associativity int, blockSize int) *MesifSimulator {
	cores := []*core.Core{}
	bus := bus.NewBus()
	memory := memory.NewMemory(constants.NumCores, bus)

	for i := 0; i < constants.NumCores; i++ {
		cache := cache.NewMesifCache(i, bus, blockSize, associativity, cacheSize)
		cores = append(cores, core.NewCore(i, inputFilePrefix, cache))
	}

	return &MesifSimulator{BaseSimulator: simulator.NewBaseSimulator(cores, bus, memory)}
}
