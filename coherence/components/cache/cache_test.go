package cache

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chriskheng/cs4223-assignment2/coherence/testutils"
)

type cacheDsTestSet struct {
	blockSize               int
	associativity           uint32
	cacheSize               int
	numSets                 uint32
	offsetNumBits           uint32
	setIndexNumBits         uint32
	cacheTests              []cacheTagAndIndexTest
	cacheContainAndSetTests []cacheContainAndInsertTest
}

type cacheTagAndIndexTest struct {
	address         uint32
	tag             uint32
	setIndex        uint32
	absoluteIndices []uint32
}

type cacheContainAndInsertTest struct {
	address      uint32
	action       cacheAction
	evictedAddrs []uint32
	isExist      bool
}

type cacheAction int

const (
	Insert cacheAction = iota
	Access
)

var tests = []cacheDsTestSet{
	{
		blockSize:       16,
		associativity:   2,
		cacheSize:       1024,
		numSets:         32,
		offsetNumBits:   4,
		setIndexNumBits: 5,
		cacheTests: []cacheTagAndIndexTest{
			{
				address:         0xCFF,
				tag:             6,
				setIndex:        0xF,
				absoluteIndices: []uint32{0xF, 0xF + 32},
			},
			{
				address:         0xFEC,
				tag:             7,
				setIndex:        0x1E,
				absoluteIndices: []uint32{0x1E, 0x1E + 32},
			},
		},
	},
}

var testsInsert = []cacheDsTestSet{
	{
		blockSize:       16,
		associativity:   2,
		cacheSize:       1024,
		numSets:         32,
		offsetNumBits:   4,
		setIndexNumBits: 5,
		cacheContainAndSetTests: []cacheContainAndInsertTest{
			{
				address: 0xFEC,
				action:  Insert,
			},
			{
				address: 0xDE0,
				action:  Insert,
			},
			{
				address:      0x1E3,
				action:       Insert,
				evictedAddrs: []uint32{0xFEC},
			},
			{
				address:      0x3E4,
				action:       Insert,
				evictedAddrs: []uint32{0xDE0},
			},
		},
	},
	{
		blockSize:       16,
		associativity:   2,
		cacheSize:       1024,
		numSets:         32,
		offsetNumBits:   4,
		setIndexNumBits: 5,
		cacheContainAndSetTests: []cacheContainAndInsertTest{
			{
				address: 0xFEC,
				action:  Insert,
			},
			{
				address: 0xDE0,
				action:  Insert,
			},
			{
				address: 0xFEC,
				action:  Access,
				isExist: true,
			},
			{
				address:      0x1E3,
				action:       Insert,
				evictedAddrs: []uint32{0xDE0},
			},
			{
				address:      0x3E4,
				action:       Insert,
				evictedAddrs: []uint32{0xFEC},
			},
		},
	},
}

func Test(t *testing.T) {
	for _, test := range tests {
		cacheDs := NewCacheDs(test.blockSize, int(test.associativity), test.cacheSize)
		if cacheDs.offsetNumBits != test.offsetNumBits {
			t.Fatalf("Incorrect offsetNumBits: expected %d, got %d", test.offsetNumBits, cacheDs.offsetNumBits)
		}

		if cacheDs.setIndexNumBits != test.setIndexNumBits {
			t.Fatalf("Incorrect indexNumBits: expected %d, got %d", test.setIndexNumBits, cacheDs.setIndexNumBits)
		}

		if cacheDs.associativity != test.associativity {
			t.Fatalf("Incorrect associativity: expected %d, got %d", test.associativity, cacheDs.associativity)
		}

		if cacheDs.numSets != test.numSets {
			t.Fatalf("Incorrect numBlocksPerSet: expected %d, got %d", test.numSets, cacheDs.numSets)
		}

		for _, subtest := range test.cacheTests {
			tag1 := cacheDs.GetTag(subtest.address)
			if tag1 != subtest.tag {
				expected := fmt.Sprintf("%x", subtest.tag)
				got := fmt.Sprintf("%x", tag1)
				t.Fatalf("Incorrect tag: expected %s, got %s", expected, got)
			}

			setIndex := cacheDs.GetCacheSetIndex(subtest.address)
			if setIndex != subtest.setIndex {
				expected := fmt.Sprintf("%x", subtest.setIndex)
				got := fmt.Sprintf("%x", setIndex)
				t.Fatalf(testutils.GetErrorString("setIndex", expected, got))
			}

			for i := 0; i < int(test.associativity); i++ {
				absoluteIndex := cacheDs.getAbsoluteIndex(subtest.address, i)
				if absoluteIndex != subtest.absoluteIndices[i] {
					identifier := fmt.Sprintf("absoluteIndex%d", i)
					expected := strconv.Itoa(int(subtest.absoluteIndices[i]))
					got := strconv.Itoa(int(absoluteIndex))
					t.Fatalf(testutils.GetErrorString(identifier, expected, got))
				}
			}
		}
	}
}

func TestInsert(t *testing.T) {
	for _, test := range testsInsert {
		cacheDs := NewCacheDs(test.blockSize, int(test.associativity), test.cacheSize)

		for _, subtest := range test.cacheContainAndSetTests {
			switch subtest.action {
			case Insert:
				isEvicted, address := cacheDs.Insert(subtest.address)
				if !cacheDs.Contain(subtest.address) {
					expected := strconv.FormatBool(true)
					got := strconv.FormatBool(false)
					t.Fatalf(testutils.GetErrorString("boolean (contain)", expected, got))
				}

				if len(subtest.evictedAddrs) == 0 && isEvicted {
					identifier := fmt.Sprintf("evicted (insert %x)", subtest.address)
					expected := strconv.FormatBool(false)
					got := strconv.FormatBool(isEvicted)
					t.Fatalf(testutils.GetErrorString(identifier, expected, got))
				} else if len(subtest.evictedAddrs) != 0 {
					if !isEvicted {
						identifier := fmt.Sprintf("evicted (insert %x)", subtest.address)
						expected := strconv.FormatBool(true)
						got := strconv.FormatBool(isEvicted)
						t.Fatalf(testutils.GetErrorString(identifier, expected, got))
					}

					if address != subtest.evictedAddrs[0] {
						expected := fmt.Sprintf("%x", subtest.evictedAddrs[0])
						got := fmt.Sprintf("%x", address)
						t.Fatalf(testutils.GetErrorString("evicted address", expected, got))
					}
				}
			case Access:
				isExist := cacheDs.Access(subtest.address)
				if isExist != subtest.isExist {
					expected := strconv.FormatBool(subtest.isExist)
					got := strconv.FormatBool(isExist)
					t.Fatalf(testutils.GetErrorString("isExist (access)", expected, got))
				}
			}
		}
	}
}
