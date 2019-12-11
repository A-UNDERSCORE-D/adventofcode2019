package intcode2

import (
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func NewFromString(code string) *IVM {
	splitCode := util.GetInts(strings.Split(code, ","))
	mem := make(map[int]int, len(splitCode))
	for i, data := range splitCode {
		mem[i] = data
	}
	return &IVM{
		Memory: mem,
		Input:  make(chan int, 20),
		Output: make(chan int),
	}
}

type IVM struct {
	Memory       map[int]int
	IP           int
	RelativeBase int
	Input        chan int
	Output       chan int

	// debug stuff
}

// TODO: <dsal> modes := make([]int, 0, 3)      Then pass in modes[:0] each time instead of making it in the function.

func (ivm *IVM) DecodeOpAt(idx int) (int, []int) {
	rawOp := ivm.Memory[idx]
	if rawOp <= 99 { // Fast path for operations without modes
		return rawOp, nil
	}
	argmodes := make([]int, 0, 3)
	// "decimal right shift" the argument by division, grab the right side with modulo
	op := rawOp % 100
	rawArgModes := rawOp / 100
	for rawArgModes > 9 {
		argmodes = append(argmodes, rawArgModes%10)
		rawArgModes /= 10
	}
	argmodes = append(argmodes, rawArgModes)
	return op, argmodes
}

func (ivm *IVM) Execute() {
	for {
		opcode, argModes := ivm.DecodeOpAt(ivm.IP)
		if opcode == Halt {
			break
		}

		ivm.IP += ivm.executeOp(opcode, argModes)
		ivm.IP++
	}
	close(ivm.Input)
	close(ivm.Output)
}

func (ivm *IVM) getArgs(count int, argModes []int, resolve ...bool) []int {
	if len(resolve) != count {
		panic("incorrect number of resolve targets for getArgs")
	}
	args := make([]int, 0, 3)
	for pos := 0; pos < count; pos++ {
		var idx int
		if pos < len(argModes) {
			idx = argModes[pos]
		}
		var ptr int
		switch idx {
		case ModePosition:
			ptr = ivm.Memory[ivm.IP+pos+1]
		case ModeImmediate:
			args = append(args, ivm.Memory[ivm.IP+pos+1])
			continue
		case ModeRelative:
			ptr = ivm.RelativeBase + ivm.Memory[ivm.IP+pos+1]
		}
		if resolve[pos] {
			args = append(args, ivm.Memory[ptr])
		} else {
			args = append(args, ptr)
		}

	}
	return args
}

func (ivm *IVM) executeOp(opcode int, argModes []int) int {
	switch opcode {
	case Add:
		args := ivm.getArgs(3, argModes, true, true, false)
		a, b, target := args[0], args[1], args[2]
		ivm.Memory[target] = a + b
		return 3
	case Mul:
		args := ivm.getArgs(3, argModes, true, true, false)
		a, b, target := args[0], args[1], args[2]
		ivm.Memory[target] = a * b
		return 3
	case Input:
		target := ivm.getArgs(1, argModes, false)[0]
		ivm.Memory[target] = <-ivm.Input
		return 1
	case Output:
		data := ivm.getArgs(1, argModes, true)[0]
		ivm.Output <- data
		return 1
	case JumpIfTrue:
		args := ivm.getArgs(2, argModes, true, true)
		toCheck, target := args[0], args[1]
		if toCheck != 0 {
			ivm.IP = target
			return -1 // -1 to deal with the ivm.IP++ in the calling code
		}
		return 2
	case JumpIfFalse:
		args := ivm.getArgs(2, argModes, true, true)
		toCheck, target := args[0], args[1]
		if toCheck == 0 {
			ivm.IP = target
			return -1 // -1 to deal with the ivm.IP++ in the calling code
		}
		return 2
	case LessThan:
		args := ivm.getArgs(3, argModes, true, true, false)
		a, b, targetPtr := args[0], args[1], args[2]
		if a < b {
			ivm.Memory[targetPtr] = 1
		} else {
			ivm.Memory[targetPtr] = 0
		}
		return 3
	case Equals:
		args := ivm.getArgs(3, argModes, true, true, false)
		a, b, targetPtr := args[0], args[1], args[2]
		if a == b {
			ivm.Memory[targetPtr] = 1
		} else {
			ivm.Memory[targetPtr] = 0
		}
		return 3
	case RelBaseOffset:
		toSet := ivm.getArgs(1, argModes, true)[0]
		ivm.RelativeBase += toSet
		return 1
	case Halt:
		panic("halt sent to executeOP when it should be handled by calling code")
	default:
		panic("unknown opcode")
	}
}
