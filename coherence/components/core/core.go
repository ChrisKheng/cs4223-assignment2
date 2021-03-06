/*
Package core implements a Core structure that simulates a core in a multi-core processor.
*/
package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/stats"
	"github.com/chriskheng/cs4223-assignment2/coherence/utils"
)

type Core struct {
	cache   cache.CacheController
	reader  *bufio.Reader
	index   int
	state   CoreState
	counter int
	stats   CoreStats
}

type CoreStats struct {
	NumComputeCycles int
	NumLoads         int
	NumStores        int
	NumIdleCycles    int
}

type CoreState int

const (
	Ready CoreState = iota
	ComputeState
	MemoryState
	Done
)

func (s CoreState) String() string {
	return [...]string{"Ready", "Compute", "Memory", "Done"}[s]
}

const (
	LoadOp   = "0"
	StoreOp  = "1"
	OthersOp = "2"
)

func NewCore(index int, inputFilePrefix string, cache cache.CacheController) *Core {
	f, err := os.Open(fmt.Sprintf("%s_%d.data", inputFilePrefix, index))
	utils.Check(err)

	reader := bufio.NewReader(f)
	return &Core{cache: cache, reader: reader, index: index, state: Ready}
}

func (core *Core) Execute() {
	if core.state == Done {
		return
	}

	if core.state == ComputeState {
		core.counter--
		core.stats.NumComputeCycles++
		if core.counter == 0 {
			core.state = Ready
		}
	} else if core.state == MemoryState {
		core.stats.NumIdleCycles++
	} else {
		line, err := core.reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			core.state = Done
			return
		}

		inst, err := parseInstruction(line)
		utils.Check(err)

		if inst.iType == othersOp {
			cycles := inst.value

			if cycles > 1 {
				core.counter = int(cycles) - 1
				core.state = ComputeState
			}
			core.stats.NumComputeCycles++
		} else if inst.iType == loadOp {
			core.cache.RequestRead(inst.value, core.OnRequestComplete)
			core.state = MemoryState
			core.stats.NumLoads++
		} else if inst.iType == storeOp {
			core.cache.RequestWrite(inst.value, core.OnRequestComplete)
			core.state = MemoryState
			core.stats.NumStores++
		} else {
			panic(errors.New("unknown operation type"))
		}
	}

	// fmt.Printf("Core %d is in %d state\n", core.index, core.state)

	core.cache.Execute()
}

func (core *Core) GetStatistics() stats.Stats {
	cacheControllerStats := core.cache.GetStats()
	return stats.Stats{
		NumComputeCycles:         core.stats.NumComputeCycles,
		NumLoads:                 core.stats.NumLoads,
		NumStores:                core.stats.NumStores,
		NumIdleCycles:            core.stats.NumIdleCycles,
		NumAccessesToPrivateData: cacheControllerStats.NumAccessesToPrivateData,
		NumAccessesToSharedData:  cacheControllerStats.NumAccessesToSharedData,
		NumCacheMisses:           cacheControllerStats.NumCacheMisses,
		NumCacheAccesses:         cacheControllerStats.NumCacheAccesses,
	}
}

func (core *Core) IsDone() bool {
	return core.state == Done
}

func (core *Core) OnRequestComplete() {
	if core.state != MemoryState {
		panic("onRequestComplete should only be called when the core is in the memory state")
	}
	core.state = Ready
}
