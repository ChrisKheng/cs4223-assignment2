package simulator

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/bus"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/components/memory"
	"github.com/chriskheng/cs4223-assignment2/coherence/stats"
)

type BaseSimulator struct {
	Cores  []core.Core
	Bus    bus.Bus
	Memory memory.Memory
}

func (s *BaseSimulator) Run() {
	for !s.isAllCoresDone() {
		for i := 0; i < len(s.Cores); i++ {
			s.Cores[i].Execute()
		}

		s.Memory.Execute()
		s.Bus.Execute()
	}

	coreStats := []core.CoreStats{}
	for i := range s.Cores {
		coreStats = append(coreStats, s.Cores[i].GetStatistics())
	}
	stats.PrintStatistics(coreStats)
}

func (s *BaseSimulator) isAllCoresDone() bool {
	for i := range s.Cores {
		if !s.Cores[i].IsDone() {
			return false
		}
	}
	return true
}
