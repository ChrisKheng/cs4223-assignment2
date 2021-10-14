package mesi

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type MesiCache struct {
	*cache.BaseCache
}

func (c *MesiCache) Execute() {
	fmt.Println("Hello from MesiCache")
}

func (c *MesiCache) RequestRead() {

}

func (c *MesiCache) RequestWrite() {

}
