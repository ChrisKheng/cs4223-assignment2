package main

import (
	"fmt"
	"os"

	"github.com/chriskheng/cs4223-assignment2/coherence/dragon"
	"github.com/chriskheng/cs4223-assignment2/coherence/mesi"
	"github.com/chriskheng/cs4223-assignment2/coherence/mesif"
	"github.com/chriskheng/cs4223-assignment2/coherence/parser"
	"github.com/chriskheng/cs4223-assignment2/coherence/simulator"
)

func main() {
	inputParser := parser.InputParser{}
	err := inputParser.Parse()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "")
		inputParser.PrintUsage()
		return
	}

	var sim simulator.Simulator
	if inputParser.Protocol == parser.Mesi {
		sim = mesi.NewMesiSimulator(inputParser.InputFileName, inputParser.CacheSize, inputParser.CacheAssociativity, inputParser.CacheBlockSize)
	} else if inputParser.Protocol == parser.Dragon {
		sim = dragon.NewDragonSimulator(inputParser.InputFileName, inputParser.CacheSize, inputParser.CacheAssociativity, inputParser.CacheBlockSize)
	} else {
		sim = mesif.NewMesifSimulator(inputParser.InputFileName, inputParser.CacheSize, inputParser.CacheAssociativity, inputParser.CacheBlockSize)
	}

	sim.Run()
}
