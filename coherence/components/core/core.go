package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/utils"
)

type Core struct {
	cache   cache.Cache
	reader  *bufio.Reader
	index   int
	state   CoreState
	counter int
	stats   CoreStats
}

type CoreStats struct {
	NumComputeCycles int
	NumLoadStores    int
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

func NewCore(index int, inputFilePrefix string, cache cache.Cache) Core {
	f, err := os.Open(fmt.Sprintf("%s_%d.data", inputFilePrefix, index))
	utils.Check(err)

	reader := bufio.NewReader(f)
	return Core{cache: cache, reader: reader, index: index, state: Ready}
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
			cycles, err := strconv.ParseInt(inst.value, 0, 64)
			if err != nil {
				panic(err)
			}

			if cycles > 1 {
				core.counter = int(cycles) - 1
				core.state = ComputeState
			}
			core.stats.NumComputeCycles++
		} else if inst.iType == loadOp {
			// TODO: Call relevant method of cache
			core.stats.NumLoadStores++
		} else if inst.iType == storeOp {
			// TODO: Call relevant method of cache
			core.stats.NumLoadStores++
		} else {
			panic(errors.New("unknown operation type"))
		}
	}

	core.cache.Execute()
}

func (core *Core) GetStatistics() CoreStats {
	return core.stats
}

func (core *Core) IsDone() bool {
	return core.state == Done
}

func (core *Core) OnRequestComplete() {

}
