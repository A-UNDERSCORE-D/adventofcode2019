package intcode

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

// INTCODE opcodes
const (
	Add           = 1
	Mul           = 2
	Input         = 3
	Output        = 4
	JumpIfTrue    = 5 // essentially JNZ
	JumpIfFalse   = 6 // essentially JEZ
	LessThan      = 7
	Equals        = 8
	RelBaseOffset = 9

	Halt = 99
)

type opInfo struct {
	name     string
	argCount int
}

var ops = map[int]opInfo{
	Add:           {"add", 3},
	Mul:           {"mul", 3},
	Input:         {"ipt", 1},
	Output:        {"out", 1},
	JumpIfTrue:    {"jnz", 2},
	JumpIfFalse:   {"jez", 2},
	LessThan:      {"lt ", 3},
	Equals:        {"eq ", 3},
	RelBaseOffset: {"rel", 1},
	Halt:          {"hlt", 0},
}

func op2Str(op int) string {
	info, known := ops[op]
	if known {
		return info.name
	}
	return "???"
}

const (
	ModePosition  = 0
	ModeImmediate = 1
	ModeRelative  = 2
)

func New(code string) *IVM {
	return NewWithMem(util.GetInts(strings.Split(code, ",")))
}
func NewWithMem(code []int) *IVM {
	mem := map[int]int{}
	for i, instruction := range code {
		mem[i] = instruction
	}
	return &IVM{
		Memory: mem,
		Input:  make(chan int, 20),
		Output: make(chan int),
	}
}

func NewWithDebug(code string, asm []string, regInfo map[int]string) *IVM {
	i := New(code)
	i.Asm = asm
	i.RegInfo = regInfo
	return i
}

type IVM struct {
	Asm        []string
	Memory     map[int]int
	IP         int
	RelBase    int
	Input      chan int
	Output     chan int
	OutputData []int
	Debug      bool
	RegInfo    map[int]string
}

func (i *IVM) getMemoryAt(address int) int {
	if address < 0 {
		panic("negative memory address")
	}
	mem, ok := i.Memory[address]
	if !ok {
		// if you access an address that doesnt exist, go creates that address in the map, this is fine in most cases
		// but causes an issue when, say, outputting your own code by iterating your memory
		delete(i.Memory, address)
	}
	return mem
}

func (i *IVM) setMemoryAt(address, data int) {
	if address < 0 {
		panic("negative memory address")
	}
	i.Memory[address] = data
}

func getOp(op int) (int, []int) {
	if op < 0 {
		return op, nil
	}
	strOP := fmt.Sprint(op)
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
			i.Print()
		}
		switch op {
		case Add:
			args := i.getArgs(2, argModes) // set the last to immed so we can access memory ourselves
			ptr := i.getMemoryAt(i.IP+3)
			if len(argModes) > 2 && argModes[2] == ModeRelative {
				ptr += i.RelBase
			}

			i.setMemoryAt(ptr, args[0]+args[1])
			i.IP += 3
		case Mul:
			args := i.getArgs(2, argModes) // set the last to immed so we can access memory ourselves
			ptr := i.getMemoryAt(i.IP+3)
			if len(argModes) > 2 && argModes[2] == ModeRelative {
				ptr += i.RelBase
			}
			i.setMemoryAt(ptr, args[0]*args[1])
			i.IP += 3
		case Input:
			ptr := i.getArgs(1, []int{ModeImmediate})[0]
			if len(argModes) > 0 && argModes[0] == ModeRelative {
				ptr += i.RelBase
			}

			i.setMemoryAt(ptr, <-i.Input)
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
			args := i.getArgs(2, argModes)
			ptr := i.Memory[i.IP+3]
			if len(argModes) > 2 && argModes[2] == ModeRelative {
				ptr += i.RelBase
			}
			if args[0] < args[1] {
				i.setMemoryAt(ptr, 1)
			} else {
				i.setMemoryAt(ptr, 0)
			}
			i.IP += 3

		case Equals:
			args := i.getArgs(2, argModes)
			ptr := i.Memory[i.IP+3]
			if len(argModes) > 2 && argModes[2] == ModeRelative {
				ptr += i.RelBase
			}
			if args[0] == args[1] {
				i.setMemoryAt(ptr, 1)
			} else {
				i.setMemoryAt(ptr, 0)
			}
			i.IP += 3

		case RelBaseOffset:
			args := i.getArgs(1, argModes)
			i.RelBase += args[0]
			i.IP += 1
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

func idxOrZero(idx int, slice []int) int {
	if idx < len(slice) {
		return slice[idx]
	}
	return ModePosition
}

func (i *IVM) getArgs(count int, modeSet []int) (args []int) {
	return i.getArgsDirect(count, i.IP, i.RelBase, modeSet)
}

func (i *IVM) getArgsDirect(count, offset, relBase int, modeSet []int) (args []int) {
	for pos := 0; pos < count; pos++ {
		switch idxOrZero(pos, modeSet) {
		case ModePosition:
			args = append(args, i.Memory[i.Memory[offset+pos+1]])
		case ModeImmediate:
			args = append(args, i.Memory[offset+pos+1])
		case ModeRelative:
			ptr := relBase + i.Memory[offset+pos+1]
			args = append(args, i.Memory[ptr])
		}
	}
	return
}

func idxOrEmptyString(idx int, slice []string) string {
	if len(slice) > idx {
		return slice[idx]
	}
	return ""
}

func getArgSpacing(in map[int]string) int {
	out := 0
	for _, v := range in {
		out = util.Max(out, len(v))
	}
	return out
}

func (i *IVM) memAddrs2Ints() []int {
	var out []int
	for k := range i.Memory {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

func (i *IVM) Print() {
	insCnt := 0
	maxArgLen := getArgSpacing(i.RegInfo) + 7

	if maxArgLen == 0 {
		maxArgLen = 10
	}
	strArg := fmt.Sprintf("%%-%ds ", maxArgLen)
	fmt.Print("┌", strings.Repeat("─", 198), "┐", "\n")
	fmt.Printf("│  IP  │ op (raw) │ args (rel: %03d) %s │ Assembly (if debug)\n", i.RelBase, strings.Repeat(" ", (maxArgLen*3)-10))
	for memAddrIdx, memaddrs := 0, i.memAddrs2Ints(); memAddrIdx < len(memaddrs); memAddrIdx++ {
		idx := memaddrs[memAddrIdx]

		out := strings.Builder{}
		op, argModes := getOp(i.Memory[idx])

		if idx == i.IP {
			out.WriteString("--> ")
		} else {
			out.WriteString("    ")
		}
		if name, exists := i.RegInfo[idx]; exists {
			out.WriteString(fmt.Sprintf("%03d %-26s: %02d", idx, name, i.Memory[idx]))
			fmt.Println(out.String())
			continue
		}
		out.WriteString(fmt.Sprintf("%03d", idx))
		out.WriteString(": ")

		info, known := ops[op]
		if known {
			out.WriteString(info.name)
			asm := idxOrEmptyString(insCnt, i.Asm)
			out.WriteString(fmt.Sprintf(" (%02d) │ ", op))

			if info.argCount > 0 {
				if len(i.Memory) < idx+info.argCount {
					out.WriteString(fmt.Sprintf("% 7d", op))
					goto end
				}

				memRng := []int{}
				for memIDX := idx + 1; memIDX < idx+1+info.argCount; memIDX++ {
					memRng = append(memRng, i.getMemoryAt(memIDX))
				}

				for argCnt, a := range memRng {
					switch idxOrZero(argCnt, argModes) {
					case ModeImmediate:
						out.WriteString(fmt.Sprintf(strArg, fmt.Sprintf("%d", a)))
					case ModeRelative:
						a = a + i.RelBase
						if rName, exists := i.RegInfo[a]; exists {
							out.WriteString(fmt.Sprintf(strArg, fmt.Sprintf("%s (%03d)", rName, i.getMemoryAt(a))))
						} else {
							out.WriteString(fmt.Sprintf(strArg, fmt.Sprintf("*%d*", a)))
						}
					case ModePosition:
						if rName, exists := i.RegInfo[a]; exists {
							out.WriteString(fmt.Sprintf(strArg, fmt.Sprintf("%s (%03d)", rName, i.getMemoryAt(a))))
						} else {
							out.WriteString(fmt.Sprintf(strArg, fmt.Sprintf("[%d]", a)))
						}
					}
					out.WriteString("│ ")
				}
				memAddrIdx += info.argCount
			}
			if info.argCount < 3 {
				out.WriteString(strings.Repeat(" ", (3-info.argCount)*(maxArgLen+3)-2))
				out.WriteString("│ ")
			}
			out.WriteString(fmt.Sprintf("%s", asm))
			insCnt++
		} else {
			out.WriteString("???   ")
			out.WriteString(fmt.Sprintf(" %8d", op))
		}
	end:
		fmt.Println(out.String())
	}
}

func (i *IVM) HandleOut() {
	go func() {
		for x := range i.Output {
			i.OutputData = append(i.OutputData, x)
		}
	}()
}
