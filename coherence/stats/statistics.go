/*
Package stats implements methods for printing formatted program statistics.
*/
package stats

import (
	"fmt"
	"os"
	"time"
)

type Stats struct {
	NumComputeCycles         int
	NumLoads                 int
	NumStores                int
	NumIdleCycles            int
	NumAccessesToPrivateData int
	NumAccessesToSharedData  int
	NumCacheMisses           int
	NumCacheAccesses         int
	NumCacheUpdates          int
}

type OtherStats struct {
	DataTrafficOnBus int // In Bytes
	NumInvalidations int
	NumUpdates       int
}

func PrintStatistics(duration time.Duration, stats []Stats, otherStats OtherStats) {
	fmt.Printf("Overall stats:\n")
	fmt.Printf("Total time taken: %d ms\n", duration.Milliseconds())
	fmt.Printf("Total Cycles: %d\n", getMaxCycles(stats))

	fmt.Printf("======================================================\n")
	fmt.Printf("Bus stats:\n")
	fmt.Printf("Total data traffic on bus: %d bytes\n", otherStats.DataTrafficOnBus)
	fmt.Printf("Total invalidations: %d\n", otherStats.NumInvalidations)
	fmt.Printf("Total updates: %d\n", otherStats.NumUpdates)

	for i := range stats {
		fmt.Printf("======================================================\n")
		fmt.Printf("Core %d:\n", i)
		fmt.Printf("Execution cycles: %d\n", getExecutionCycles(stats[i]))
		fmt.Printf("Compute cycles: %d\n", stats[i].NumComputeCycles)
		fmt.Printf("Num loads: %d\n", stats[i].NumLoads)
		fmt.Printf("Num stores: %d\n", stats[i].NumStores)
		fmt.Printf("Idle cycles: %d\n", stats[i].NumIdleCycles)
		fmt.Printf("Data cache miss rate: %.3f\n", getCacheMissRate(stats[i]))
		fmt.Printf("Num accesses to private data: %d\n", stats[i].NumAccessesToPrivateData)
		fmt.Printf("Num accesses to shared data: %d\n", stats[i].NumAccessesToSharedData)
	}
}

func PrintStatisticsCsv(duration time.Duration, stats []Stats, otherStats OtherStats) {
	fmt.Fprintf(os.Stderr, "%d,%d\n", duration.Milliseconds(), getMaxCycles(stats))
	fmt.Fprintf(os.Stderr, "%d,%d,%d\n", otherStats.DataTrafficOnBus, otherStats.NumInvalidations, otherStats.NumUpdates)

	for i := range stats {
		fmt.Fprintf(os.Stderr, "%d,%d,%d,%d,%d,%d,%.3f,%d,%d\n",
			i,
			getExecutionCycles(stats[i]),
			stats[i].NumComputeCycles,
			stats[i].NumLoads,
			stats[i].NumStores,
			stats[i].NumIdleCycles,
			getCacheMissRate(stats[i]),
			stats[i].NumAccessesToPrivateData,
			stats[i].NumAccessesToSharedData,
		)
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
