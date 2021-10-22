package stats

import (
	"fmt"
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/components/core"
)

func PrintStatistics(duration time.Duration, stats []core.CoreStats) {
	fmt.Printf("Total time taken: %d ms\n", duration.Milliseconds())
	fmt.Printf("Total Cycles: %d\n", getMaxCycles(stats))
	for i := range stats {
		fmt.Printf("======================================================\n")
		fmt.Printf("Core %d:\n", i)
		fmt.Printf("Execution Cycles: %d\n", getExecutionCycles(stats[i]))
		fmt.Printf("Compute Cycles: %d\n", stats[i].NumComputeCycles)
		fmt.Printf("Loads and Stores: %d\n", stats[i].NumLoadStores)
		fmt.Printf("Idle Cycles: %d\n", stats[i].NumIdleCycles)
	}
}

func getMaxCycles(stats []core.CoreStats) int {
	max := 0
	for i := range stats {
		numCycles := getExecutionCycles(stats[i])
		if numCycles > max {
			max = numCycles
		}
	}
	return max
}

func getExecutionCycles(stats core.CoreStats) int {
	return stats.NumComputeCycles + stats.NumIdleCycles
}
