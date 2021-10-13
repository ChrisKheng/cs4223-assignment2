package mesi

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

type MesiSimulator struct {
	*simulator.BaseSimulator
}

func NewMesiSimulator(numCores int, inputFilePrefix string, cacheSize int, associativity int, blockSize int) *MesiSimulator {
	cores := []core.Core{}

	for i := 0; i < numCores; i++ {
		cores = append(cores, core.NewCore(i, inputFilePrefix, &MesiCache{}))
	}

	return &MesiSimulator{&simulator.BaseSimulator{Cores: cores}}
}

func (s *MesiSimulator) Run() {
	for i := 0; i < len(s.Cores); i++ {
		s.Cores[i].Execute()
	}
}
