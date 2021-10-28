package cache

import (
	"math"
	"time"

	"github.com/chriskheng/cs4223-assignment2/coherence/constants"
)

type Cache struct {
	offsetNumBits    uint32
	setIndexNumBits  uint32
	associativity    uint32
	numSets          uint32
	blockSizeInWords uint32
	cacheArray       []CacheLine
}

type CacheLine struct {
	tag       uint32
	address   uint32
	timestamp int64 // Resolution: ns
}

// Return a new Cache struct based on the given parameters.
// All parameters given should NOT be in the form of log2(x), where x is the value of the parameter,
// i.e. the value given should be x, not log2(x).
// blockSize and cacheSize are in unit of bytes.
func NewCacheDs(blockSize, associativity, cacheSize int) Cache {
	numBlocks := cacheSize / blockSize
	numSets := uint32(numBlocks / associativity)

	return Cache{
		offsetNumBits:    uint32(math.Log2(float64(blockSize))),
		setIndexNumBits:  uint32(math.Log2(float64(numSets))),
		associativity:    uint32(associativity),
		numSets:          numSets,
		blockSizeInWords: uint32(blockSize) / constants.WordSize,
		cacheArray:       make([]CacheLine, numBlocks),
	}
}

// Return the index of the address in the underlying array if the data at the address is cached,
// otherwise return -1.
func (cacheDs *Cache) GetIndexInArray(address uint32) int {
	tag := cacheDs.GetTag(address)

	for i := 0; i < int(cacheDs.associativity); i++ {
		absoluteIndex := cacheDs.getAbsoluteIndex(address, i)
		cacheLine := &cacheDs.cacheArray[absoluteIndex]
		if cacheLine.timestamp != 0 && cacheLine.tag == tag {
			return int(absoluteIndex)
		}
	}

	return -1
}

// Return true if the cache contains the data at the address.
func (cacheDs *Cache) Contain(address uint32) bool {
	return cacheDs.GetIndexInArray(address) != -1
}

// Return true together with the evicted address if eviction occurs during the insertion.
// Also return the index of the position in the array where the cache line is inserted.
func (cacheDs *Cache) Insert(address uint32) (bool, uint32, int) {
	// If the address is already in the cache, then skip the insert
	if cacheDs.GetIndexInArray(address) != -1 {
		return false, 0, -1
	}

	index := cacheDs.getAbsoluteIndex(address, 0)
	leastTimestamp := cacheDs.cacheArray[index].timestamp

	// LRU eviction policy
	for i := 1; i < int(cacheDs.associativity); i++ {
		absoluteIndex := cacheDs.getAbsoluteIndex(address, i)
		cacheLine := &cacheDs.cacheArray[absoluteIndex]
		if cacheLine.timestamp < leastTimestamp {
			leastTimestamp = cacheLine.timestamp
			index = absoluteIndex
		}
	}

	isToBeEvicted := cacheDs.cacheArray[index].timestamp != 0
	evictedAddress := cacheDs.cacheArray[index].address

	cacheDs.cacheArray[index] = CacheLine{
		tag:       cacheDs.GetTag(address),
		address:   address,
		timestamp: time.Now().UnixNano(),
	}

	// This is needed to add some delay between two consecutive accesses, otherwise the test cases may fail some time
	// as the time between two consecutive accesses would become the same
	time.Sleep(time.Microsecond * 1)

	return isToBeEvicted, evictedAddress, int(index)
}

// Access the cache line of the given address and update the cache line timestamp.
// Return true if the data at the address is cached, otherwise false.
func (cacheDs *Cache) Access(address uint32) bool {
	index := cacheDs.GetIndexInArray(address)
	if index == -1 {
		return false
	}
	cacheDs.cacheArray[index].timestamp = time.Now().UnixNano()
	return true
}

// Remove a cache line from the cache.
func (cacheDs *Cache) Evict(address uint32) {
	index := cacheDs.GetIndexInArray(address)
	cacheDs.cacheArray[index].timestamp = 0
}

// Return the tag of the given address.
func (cacheDs *Cache) GetTag(address uint32) uint32 {
	return address >> (cacheDs.setIndexNumBits + cacheDs.offsetNumBits)
}

// Return the set index of the given address.
func (cacheDs *Cache) GetCacheSetIndex(address uint32) uint32 {
	return (address >> (cacheDs.offsetNumBits)) & ((1 << cacheDs.setIndexNumBits) - 1)
}

// Return the index of the given address in the underlying array based on the given round parameter.
// It computes the index by first getting the set index of the given address, and add round * number of cache sets.
func (cacheDs *Cache) getAbsoluteIndex(address uint32, round int) uint32 {
	normalizedIndex := cacheDs.GetCacheSetIndex(address)
	return normalizedIndex + uint32(round)*cacheDs.numSets
}
