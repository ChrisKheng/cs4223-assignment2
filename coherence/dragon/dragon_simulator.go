package dragon

import (
	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

type DragonSimulator struct {
	*simulator.BaseSimulator
}

func NewDragonSimulator(numCores int, inputFilePrefix string, cacheSize int, associativity int, blockSize int) *DragonSimulator {
	cores := []core.Core{}

	for i := 0; i < numCores; i++ {
		cores = append(cores, core.NewCore(i, inputFilePrefix, NewDragonCache(blockSize, associativity, cacheSize)))
	}

	return &DragonSimulator{&simulator.BaseSimulator{Cores: cores}}
}
