package asm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

// asm Operations
const (
	OpAdd         = intcode.Add
	OpMul         = intcode.Mul
	OpInput       = intcode.Input
	OpOutput      = intcode.Output
	OpJumpIfTrue  = intcode.JumpIfTrue  // essentially JNZ
	OpJumpIfFalse = intcode.JumpIfFalse // essentially JEZ
	OpLessThan    = intcode.LessThan
	OpEquals      = intcode.Equals
	OpHalt        = intcode.Halt
	OpInit        = -1
)

var opArgCounts = map[int]int{
	OpAdd:         3,
	OpMul:         3,
	OpInput:       1,
	OpOutput:      1,
	OpJumpIfTrue:  2,
	OpJumpIfFalse: 2,
	OpLessThan:    3,
	OpEquals:      3,
	OpInit:        2,
}

var ops = map[int]string{
	OpAdd:         "add",
	OpMul:         "mul",
	OpInput:       "ipt",
	OpOutput:      "out",
	OpJumpIfTrue:  "jnz",
	OpJumpIfFalse: "jez",
	OpLessThan:    "lt",
	OpEquals:      "eq",
	OpHalt:        "halt",
	OpInit:        "init",
}

func op2Str(op int) string {
	if res, ok := ops[op]; ok {
		return res
	}
	return "???"
}

func str2op(op string) (int, bool) {
	for k, v := range ops {
		if v == op {
			return k, true
		}
	}
	return -1, false
}

const (
	argImmediate = iota
	argPosition
	argRegister
	argRegisterPtr // literal location for register, ie, a pointer to that register in literal form
	argCodePos
	argCodePosIns
)

type arg struct {
	orig string
	typ  int
	id   string
}

type argSlice []arg

func (a argSlice) argModes() (out []string) {
	out = make([]string, len(a))
	for i, arg := range a {
		switch arg.typ {
		case argImmediate, argCodePos, argCodePosIns, argRegisterPtr:
			out[i] = "1"
		case argRegister, argPosition:
			out[i] = "0"
		}
	}
	for i := len(out)/2 - 1; i >= 0; i-- {
		opp := len(out) - 1 - i
		out[i], out[opp] = out[opp], out[i]
	}
	return out
}

func newArg(in string) arg {
	out := arg{orig: in}
	switch in[0] {
	case 'r':
		out.typ = argRegister
		in = in[1:]
	case '$':
		out.typ = argPosition
		in = in[1:]
	case '.':
		out.typ = argCodePos
		if len(in) == 1 {
			in = "0"
			break
		}
		in = in[1:]
	case '!':
		out.typ = argCodePosIns
		if len(in) == 1 {
			in = "0"
			break
		}
		in = in[1:]
	case '&':
		out.typ = argRegisterPtr
		in = in[1:]
	default:
		out.typ = argImmediate
	}
	out.id = in
	return out
}

type Token struct {
	orig   string
	opCode int
	args   []arg
}

var replRe = regexp.MustCompile("(.*?)#.*")

func splitAndClean(in string) (out []string) {
	for _, s := range strings.Split(strings.Replace(in, "\n", ";", -1), ";") {
		s = replRe.ReplaceAllString(s, "$1")
		if len(s) > 0 {
			trimmed := strings.TrimSpace(s)
			out = append(out, trimmed)
		}
	}
	return
}

func cleanSplit(in []string) (out []string) {
	for _, s := range in {
		s := strings.TrimSpace(s)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return
}

func getToken(in string) (Token, error) {
	split := strings.Split(in, " ")
	opStr, argStrs := strings.ToLower(split[0]), cleanSplit(split[1:])
	op, ok := str2op(opStr)
	if !ok {
		return Token{}, fmt.Errorf("unknown op %s", opStr)
	}
	if len(argStrs) != opArgCounts[op] {
		return Token{}, fmt.Errorf("incorrect argument count. got %d, want %d (from %s)", len(argStrs), opArgCounts[op], in)
	}
	var args []arg
	for _, a := range argStrs {
		args = append(args, newArg(a))
	}

	return Token{
		orig:   in,
		opCode: op,
		args:   args,
	}, nil
}

func Tokenise(in string) []Token {
	split := splitAndClean(in)
	var out []Token
	for _, s := range split {
		tok, err := getToken(s)
		if err != nil {
			panic(err)
		}
		out = append(out, tok)
	}
	return out
}

type register struct {
	codeLoc int
	content int
}

func calculateLength(tokens []Token, noAutoHalt bool) (int, []int) {
	l := 0
	var out []int
	var seenHalt bool
	for _, token := range tokens {
		if token.opCode < 0 {
			continue
		}

		l += 1
		out = append(out, l)
		l += len(token.args)
		if token.opCode == OpHalt {
			seenHalt = true
		}
	}
	if !noAutoHalt && !seenHalt {
		l++
	}
	return l, out
}

func Assemble(in []Token, noAutoHalt bool) (asm string, regInfo map[int]string, err error) {
	// TODO: enough of this single function holds all the state shit. Make. A. Struct.
	registers := map[string]register{}
	var outOpcodes []string
	codeLen, insStarts := calculateLength(in, noAutoHalt)
	curPos := -1
	curIns := -1
	var getReg func(a arg, def int) int
	getReg = func(a arg, def int) int {
		switch a.typ {
		case argCodePos:
			return curPos - 1 + util.GetInt(a.id)
		case argCodePosIns:
			diff := util.GetInt(a.id)
			return insStarts[curIns+diff] - 1
		case argRegisterPtr:
			return getReg(arg{
				orig: a.orig,
				typ:  argRegister,
				id:   a.id[1:],
			}, def)
		case argRegister:
			if a.id != "0" {
				r := registers[a.id]

				if r.codeLoc == 0 {
					r.codeLoc = codeLen
					r.content = def
					registers[a.id] = r

					codeLen++
				}
				return r.codeLoc
			}

			fallthrough
		default:
			return util.GetInt(a.id)
		}
	}

	defer func() {
		// if panik := recover(); panik != nil {
		// 	err = fmt.Errorf("could not assemble, hit error at position %d: %s", curPos, panik)
		// }
	}()

	var seenHalt bool
	for _, token := range in {
		if token.opCode < 0 {
			// Special ops to control the assembler
			switch token.opCode {
			case OpInit:
				// Initialise a register with a value
				// init register number
				regStr := token.args[0]
				toSet := util.GetInt(token.args[1].orig) // we just want the int value, nothing magic
				getReg(regStr, toSet)
				continue
			}
		}

		curPos++
		curIns++
		outOpcodes = append(outOpcodes, strings.Join(argSlice(token.args).argModes(), "")+fmt.Sprintf("%02d", token.opCode))
		if token.opCode == OpHalt {
			seenHalt = true
		}
		for _, a := range token.args {
			outOpcodes = append(outOpcodes, strconv.Itoa(getReg(a, 0)))
			curPos++
		}
	}

	if !noAutoHalt && !seenHalt {
		outOpcodes = append(outOpcodes, strconv.Itoa(OpHalt))
	}

	for _, r := range registers {
		outOpcodes = append(outOpcodes, strconv.Itoa(r.content))
	}
	fmt.Println(registers)
	return strings.Join(outOpcodes, ","), makeRegInfo(registers), nil
}

func makeRegInfo(in map[string]register) map[int]string {
	out := make(map[int]string)
	for k, v := range in {
		out[v.codeLoc] = k
	}
	return out
}
