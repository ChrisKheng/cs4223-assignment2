package core

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
	"github.com/chriskheng/cs4223-assignment2/coherence/utils"
)

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

type Core struct {
	cache  cache.Cache
	reader *bufio.Reader
	index  int
}

func NewCore(index int, inputFilePrefix string, cache cache.Cache) Core {
	f, err := os.Open(fmt.Sprintf("%s_%d.data", inputFilePrefix, index))
	utils.Check(err)

	reader := bufio.NewReader(f)
	return Core{cache: cache, reader: reader}
}

func (core *Core) Execute() {
	core.cache.Execute()
}

func (core *Core) OnRequestComplete() {

}
