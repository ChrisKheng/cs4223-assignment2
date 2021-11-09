/*
Package parser implements an InputParser struct to parse user-given arguments.
*/
package parser

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/chriskheng/cs4223-assignment2/coherence/constants"
)

type CCProtocol int

const (
	Mesi CCProtocol = iota
	Mesif
	Dragon
)

type InputParser struct {
	Protocol           CCProtocol
	InputFileName      string
	CacheSize          int
	CacheAssociativity int
	CacheBlockSize     int
}

func (p *InputParser) Parse() (err error) {
	args := os.Args[1:]
	if !(len(args) == 2 || len(args) == 5) {
		return errors.New("incorrect number of arguments provided")
	}

	if err = p.parseProtocolAndBenchmark(args[0:2]); err != nil {
		return
	}

	if len(args) == 5 {
		err = p.parseCacheConfigs(args[2:])
	} else {
		p.CacheSize = 4096
		p.CacheAssociativity = 2
		p.CacheBlockSize = 32
		fmt.Printf("Using default cache config => cache size: %dB, associativity: %d, block size: %dB\n",
			p.CacheSize, p.CacheAssociativity, p.CacheBlockSize)
	}

	return
}

func (p *InputParser) parseProtocolAndBenchmark(args []string) (err error) {
	protocol, err := p.parseProtocol(args[0])
	if err != nil {
		return
	}
	p.Protocol = protocol
	p.InputFileName = args[1]
	return
}

func (p *InputParser) parseCacheConfigs(configs []string) (err error) {
	cacheSizeValue, err1 := strconv.Atoi(configs[0])
	associativityValue, err2 := strconv.Atoi(configs[1])
	blockSizeValue, err3 := strconv.Atoi(configs[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("cache_size, associativity, or block_size provided is not an integer")
	}

	p.CacheSize = cacheSizeValue
	p.CacheAssociativity = associativityValue
	p.CacheBlockSize = blockSizeValue

	if err = p.checkCacheValues(); err != nil {
		return
	}

	return
}

func (p *InputParser) parseProtocol(protocol string) (CCProtocol, error) {
	switch protocol {
	case "MESI":
		return Mesi, nil
	case "Dragon":
		return Dragon, nil
	case "MESIF":
		return Mesif, nil
	default:
		return -1, errors.New("invalid protocol")
	}
}

func (p *InputParser) checkCacheValues() error {
	if p.CacheAssociativity != 1 && p.CacheAssociativity%2 != 0 {
		return errors.New("cache associativity to be power of 2")
	}

	if p.CacheSize%2 != 0 || p.CacheBlockSize%2 != 0 {
		return errors.New("cache_size and block_size needs to be power of 2")
	}

	if p.CacheBlockSize < int(constants.WordSize) {
		return errors.New("block_size needs to be at least the word size")
	}

	if p.CacheSize%p.CacheBlockSize != 0 {
		return errors.New("cache_size needs to be divisible by block_size")
	}

	if (p.CacheSize/p.CacheBlockSize)%p.CacheAssociativity != 0 {
		return errors.New("number of cache blocks (cache_size / block_size) needs to be divisible by associativity")
	}

	return nil
}

func (p *InputParser) PrintUsage() {
	fmt.Fprintln(os.Stderr, "Usage: coherence <protocol> <input_file_prefix> [cache_size] [associativity] [block_size]")
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintln(os.Stderr, "protocol: MESI or Dragon")
	fmt.Fprintln(os.Stderr, "input_file_prefix: Prefix to the benchmark file, "+
		"e.g. ../benchmarks/blackscholes_four/blackscholes")
	fmt.Fprintln(os.Stderr, "cache_size: cache size in bytes. Must be power of 2 and divisible by block_size")
	fmt.Fprintln(os.Stderr, "associativity: associativity of the cache. Must be power of 2 and able "+
		"to divide the number of cache sets")
	fmt.Fprintln(os.Stderr, "block_size: block size in bytes. Must be power of 2 and at least the size of a word (4 bytes).")

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "You can just provide the arguments: protocol and input_file_prefix. In this case, "+
		"the default cache configuration will be used.")
	fmt.Fprintln(os.Stderr, "Default cache configuration => cache_size: 4096B, associativity: 2, block_size: 32B")
}
