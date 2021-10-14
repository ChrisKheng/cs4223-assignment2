package core

import (
	"errors"
	"strings"
)

type instructionType int

const (
	loadOp instructionType = iota
	storeOp
	othersOp
)

type instruction struct {
	iType instructionType
	value string
}

func parseInstruction(line string) (instruction, error) {
	tokens := strings.Fields(line)
	if len(tokens) != 2 {
		return instruction{}, errors.New("illegal instruction format")
	}

	iType, err := parseInstructionType(tokens[0])
	if err != nil {
		return instruction{}, err
	}

	return instruction{iType: iType, value: tokens[1]}, nil
}

func parseInstructionType(token string) (instructionType, error) {
	switch token {
	case "0":
		return loadOp, nil
	case "1":
		return storeOp, nil
	case "2":
		return othersOp, nil
	default:
		return -1, errors.New("illegal instruction type")
	}
}
