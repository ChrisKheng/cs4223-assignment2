package simulator

import (
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/stats"
)

type BaseSimulator struct {
	cores  []*core.Core
	bus    *bus.Bus
	memory *memory.Memory
}

func NewBaseSimulator(cores []*core.Core, bus *bus.Bus, memory *memory.Memory) *BaseSimulator {
	return &BaseSimulator{
		cores:  cores,
		bus:    bus,
		memory: memory,
	}
}

func (s *BaseSimulator) Run() {
	start := time.Now()
	iter := 0
	for !s.isAllCoresDone() {
		// fmt.Printf("Current iter %d\n", iter)
		for i := 0; i < len(s.cores); i++ {
			s.cores[i].Execute()
		}

		s.bus.Execute()
		s.memory.Execute()
		iter++
	}
	elapsed := time.Since(start)

	coreStats := []stats.Stats{}
	for i := range s.cores {
		coreStats = append(coreStats, s.cores[i].GetStatistics())
	}

	otherStats := stats.OtherStats{
		DataTrafficOnBus: s.bus.GetStatistics().DataTraffic,
	}

	stats.PrintStatistics(elapsed, coreStats, otherStats)
}

func (s *BaseSimulator) isAllCoresDone() bool {
	for i := range s.cores {
		if !s.cores[i].IsDone() {
			return false
		}
	}
	return true
}
