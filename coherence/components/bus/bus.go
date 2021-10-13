package bus

import "github.com/chriskheng/cs4223-assignment2/coherence/components/memory"

type Bus struct {
	memory     memory.Memory
	busClients []BusClient
}

func (b *Bus) Execute() {

}
