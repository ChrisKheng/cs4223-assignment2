package stats

import (
	"fmt"
	"time"
)

type Stats struct {
	NumComputeCycles         int
	NumLoadStores            int
	NumIdleCycles            int
	NumAccessesToPrivateData int
	NumAccessesToSharedData  int
	NumCacheMisses           int
	NumCacheAccesses         int
}

type OtherStats struct {
	DataTrafficOnBus int // In Bytes
}

func PrintStatistics(duration time.Duration, stats []Stats, otherStats OtherStats) {
	fmt.Printf("Total time taken: %d ms\n", duration.Milliseconds())
	fmt.Printf("Total Cycles: %d\n", getMaxCycles(stats))
	fmt.Printf("Total data traffic on Bus: %d bytes\n", otherStats.DataTrafficOnBus)
	for i := range stats {
		fmt.Printf("======================================================\n")
		fmt.Printf("Core %d:\n", i)
		fmt.Printf("Execution Cycles: %d\n", getExecutionCycles(stats[i]))
		fmt.Printf("Compute Cycles: %d\n", stats[i].NumComputeCycles)
		fmt.Printf("Loads and Stores: %d\n", stats[i].NumLoadStores)
		fmt.Printf("Idle Cycles: %d\n", stats[i].NumIdleCycles)
		fmt.Printf("Data Cache Miss Rate: %.3f\n", getCacheMissRate(stats[i]))
		fmt.Printf("Num Accesses to Private Data: %d\n", stats[i].NumAccessesToPrivateData)
		fmt.Printf("Num Accesses to Shared Data: %d\n", stats[i].NumAccessesToSharedData)
	}
}

func getMaxCycles(stats []Stats) int {
	max := 0
	for i := range stats {
		numCycles := getExecutionCycles(stats[i])
		if numCycles > max {
			max = numCycles
		}
	}
	return max
}

func getExecutionCycles(stats Stats) int {
	return stats.NumComputeCycles + stats.NumIdleCycles
}

func getCacheMissRate(stats Stats) float64 {
	return float64(stats.NumCacheMisses) / float64(stats.NumCacheAccesses)
}
