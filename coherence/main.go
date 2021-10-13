package main

import (
	"fmt"
	"os"
)

func main() {
	parser := Parser{}
	err := parser.Parse()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "")
		parser.PrintUsage()
	}
}
