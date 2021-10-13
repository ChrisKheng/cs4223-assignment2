package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Parser struct {
	Protocol           string
	InputFileName      string
	CacheSize          int
	CacheAssociativity int
	CacheBlockSize     int
}

func (p *Parser) Parse() error {
	args := os.Args[1:]
	if len(args) != 5 {
		return errors.New("incorrect number of arguments provided")
	}

	p.Protocol = args[0]
	p.InputFileName = args[1]
	cacheSizeValue, err1 := strconv.Atoi(args[2])
	associativityValue, err2 := strconv.Atoi(args[3])
	blockSizeValue, err3 := strconv.Atoi(args[4])

	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("cache_size, associativity, or block_size provided is not an integer")
	}

	p.CacheSize = cacheSizeValue
	p.CacheAssociativity = associativityValue
	p.CacheBlockSize = blockSizeValue

	p.checkCacheValues()

	return nil
}

func (p *Parser) checkCacheValues() error {
	if p.CacheSize%2 != 0 || p.CacheAssociativity%2 != 0 || p.CacheBlockSize%2 != 0 {
		return errors.New("cache_size, block_size, and associativity needs to be power of 2")
	}

	if p.CacheSize%p.CacheBlockSize != 0 {
		return errors.New("cache_size needs to be divisible by block_size")
	}

	if (p.CacheSize/p.CacheBlockSize)%p.CacheAssociativity != 0 {
		return errors.New("number of cache blocks (cache_size / block_size) needs to be divisible by associativity")
	}

	return nil
}

func (p *Parser) PrintUsage() {
	fmt.Fprintln(os.Stderr, "Usage: coherence protocol input_file cache_size associativity block_size")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "protocol: MESI or Dragon")
	fmt.Fprintln(os.Stderr, "input_file: Name of benchmark. Either blackscholes, bodytrack, or fluidanimate")
	fmt.Fprintln(os.Stderr, "cache_size: cache size in bytes. Must be power of 2 and divisible by block_size")
	fmt.Fprintln(os.Stderr, "associativity: associativity of the cache. Must be power of 2 and able to divide the number of cache sets")
	fmt.Fprintln(os.Stderr, "block_size: block size in bytes. Must be power of 2.")
}
