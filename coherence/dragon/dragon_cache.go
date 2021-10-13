package dragon

import (
	"fmt"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/cache"
)

type DragonCache struct {
	*cache.BaseCache
}

func (c *DragonCache) Execute() {
	fmt.Println("Hello from Dragon Cache")
}

func (c *DragonCache) RequestRead() {

}

func (c *DragonCache) RequestWrite() {

}
