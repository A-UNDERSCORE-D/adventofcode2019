package intcode

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

// INTCODE opcodes
const (
	Add         = 1
	Mul         = 2
	Input       = 3
	Output      = 4
	JumpIfTrue  = 5 // essentially JNZ
	JumpIfFalse = 6 // essentially JEZ
	LessThan    = 7
	Equals      = 8

	Halt = 99
)

func op2Str(op int) string {
	switch op {
	case Add:
		return "add"
	case Mul:
		return "mul"
	case Input:
		return "ipt"
	case Output:
		return "out"
	case JumpIfTrue:
		return "jit"
	case JumpIfFalse:
		return "jif"
	case LessThan:
		return "lt"
	case Equals:
		return "eq"
	case Halt:
		return "halt"
	default:
		return "???"
	}
}

const (
	ModePosition  = 0
	ModeImmediate = 1
)

func New(code string) *IVM {
	return NewWithMem(util.GetInts(strings.Split(code, ",")))
}
func NewWithMem(code []int) *IVM {
	return &IVM{
		Memory: code,
		Input:  make(chan int, 1),
		Output: make(chan int),
	}
}

type IVM struct {
	Memory []int
	IP     int
	Input  chan int
	Output chan int
	OutputData []int
	Debug  bool
}

func getOp(op int) (int, []int) {
	strOP := strconv.Itoa(op)
	l := len(strOP)
	if len(strOP) == 1 {
		op = util.GetInt(strOP)
		return op, nil
	} else {
		op = util.GetInt(strOP[l-2:])
	}
	modes := strOP[:l-2]
	var argmodes []int
	for i := len(modes) - 1; i >= 0; i-- {
		argmodes = append(argmodes, util.GetInt(string(modes[i])))
	}
	return op, argmodes
}

func (i *IVM) printIfDebug(msgs ...interface{}) {
	if i.Debug {
		fmt.Println(msgs...)
	}
}

func (i *IVM) Execute() {
progLoop:
	for {
		op := i.Memory[i.IP]
		op, argModes := getOp(op)
		if i.Debug {
			fmt.Printf("IP: %03d op: %s (%02d) argModes: %v\n", i.IP, op2Str(op), op, argModes)
			i.Print()
		}
		switch op {
		case Add:
			args := i.getArgs(3, setIdx(2, ModeImmediate, 0, argModes)) // set the last to immed so we can access memory ourselves
			i.Memory[args[2]] = args[0] + args[1]
			i.IP += 3
		case Mul:
			args := i.getArgs(3, setIdx(2, ModeImmediate, 0, argModes)) // set the last to immed so we can access memory ourselves
			i.Memory[args[2]] = args[0] * args[1]
			i.IP += 3
		case Input:
			ptr := i.getArgs(1, []int{ModeImmediate})[0]
			i.Memory[ptr] = <-i.Input
			i.IP += 1
		case Output:
			arg := i.getArgs(1, argModes)[0]
			i.Output <- arg
			if i.Debug {
				runtime.Gosched()
			}
			i.IP += 1

		case JumpIfTrue:
			args := i.getArgs(2, argModes)
			if args[0] != 0 {
				i.IP = args[1]
				i.IP--
			} else {
				i.IP += 2
			}
		case JumpIfFalse:
			args := i.getArgs(2, argModes)
			if args[0] == 0 {
				i.IP = args[1]
				i.IP--
			} else {
				i.IP += 2
			}
		case LessThan:
			args := i.getArgs(3, setIdx(2, ModeImmediate, 0, argModes))
			if args[0] < args[1] {
				i.Memory[args[2]] = 1
			} else {
				i.Memory[args[2]] = 0
			}
			i.IP += 3

		case Equals:
			args := i.getArgs(3, setIdx(2, ModeImmediate, 0, argModes))
			if args[0] == args[1] {
				i.Memory[args[2]] = 1
			} else {
				i.Memory[args[2]] = 0
			}
			i.IP += 3

		case Halt:
			break progLoop
		default:
			panic(fmt.Sprintf("unexpected opcode %d at position %d", op, i.IP))
		}
		i.IP++
	}
	close(i.Input)
	close(i.Output)
}

func setIdx(idx, toSet, filler int, slice []int) []int {
	if idx >= len(slice) {
		for i := len(slice) - 1; i < idx; i++ {
			slice = append(slice, filler)
		}
	}
	slice[idx] = toSet
	return slice
}

func idxOrZero(idx int, slice []int) int {
	if idx < len(slice) {
		return slice[idx]
	}
	return ModePosition
}

func (i *IVM) getArgs(count int, modeSet []int) (args []int) {
	for pos := 0; pos < count; pos++ {
		switch idxOrZero(pos, modeSet) {
		case ModePosition:
			args = append(args, i.Memory[i.Memory[i.IP+pos+1]])
		case ModeImmediate:
			args = append(args, i.Memory[i.IP+pos+1])
		}
	}
	return
}

func (i *IVM) Print() {
	for idx, v := range i.Memory {
		if idx == i.IP {
			fmt.Printf("[%02d]", v)
		} else {
			fmt.Printf(" %02d ", v)
		}
	}
	fmt.Println()
}

func (i *IVM) HandleOut() {
	go func() {
		for x := range i.Output {
			i.OutputData = append(i.OutputData, x)
		}
	}()
}
