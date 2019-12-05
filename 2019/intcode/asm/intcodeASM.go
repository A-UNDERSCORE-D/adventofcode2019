package asm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

/*
 add %reg literal savereg
 mul %reg literal savereg
 halt does what it says on the tin
*/

const

type register struct {
	id       int
	content  int
	niceName string
	codeLoc  int
}

func (r *register) String() string {
	return r.niceName
}

type token struct {
	operation int
	opLength  int
	args      []*register
}

func (t token) String() string {
	op := ""
	switch t.operation {
	case intcode.Add:
		op = "add"
	case intcode.Mul:
		op = "mul"
	case intcode.Halt:
		op = "halt"
	}
	return fmt.Sprintf("%s: %v", op, t.args)
}

type AsmState struct {
	saveReg          int
	registers        map[int]*register
	literalRegisters []*register
	tokens           []token
}

func checkLength(tocheck []string, count int) {
	if len(tocheck) < count {
		panic("Too few arguments")
	}
}

func (a *AsmState) hasLiteral(num int) int {
	for i, r := range a.literalRegisters {
		if r.content == num {
			return i
		}
	}
	return -1
}

func (a *AsmState) GetRegister(arg string) *register {
	switch arg[0] {
	case '%':
		// we were asked for a specific register
		regNum := util.GetInt(arg[1:])
		if a.registers == nil {
			a.registers = make(map[int]*register)
		}

		if a.registers[regNum] == nil {
			a.registers[regNum] = &register{id: regNum, content: 0, niceName: "register_" + arg[1:]}
		}
		return a.registers[regNum]
	case '$':
		// we were asked for an absolute register
		ptr := util.GetInt(arg[1:])
		return &register{
			id:       -1,
			content:  0,
			niceName: fmt.Sprintf("abs_%d", ptr),
			codeLoc:  ptr,
		}
	default:
		// Literal register, store the number somewhere for them to use
		numWanted := util.GetInt(arg)
		if idx := a.hasLiteral(numWanted); idx != -1 {
			return a.literalRegisters[idx]
		}
		r := &register{id: len(a.literalRegisters), niceName: fmt.Sprintf("literal_%d", len(a.literalRegisters)), content: numWanted}
		a.literalRegisters = append(a.literalRegisters, r)
		return r
	}
}

func (a *AsmState) getRegistersForArgs(args ...string) []*register {
	var out []*register
	for _, v := range args {
		out = append(out, a.GetRegister(v))
	}
	return out
}

var replRe = regexp.MustCompile("(.*)#.*")

func splitString(in string) []string {
	var out []string
	for _, s := range strings.Split(in, "\n") {
		for _, v := range strings.Split(s, ";") {
			if v == "" {
				continue
			}
			out = append(out, strings.Trim(v, " "))
		}

	}
	return out
}

func (a *AsmState) Tokenise(in string) {
	var tokens []token
	split := splitString(in)
	for _, v := range split {
		v = replRe.ReplaceAllString(v, "$1")
		inst := strings.Split(v, " ")
		if len(inst) == 0 || inst[0] == "" {
			continue
		}
		switch strings.ToLower(inst[0]) {
		case "add":
			checkLength(inst, 4)

			tokens = append(tokens, token{
				operation: intcode.Add,
				args:      a.getRegistersForArgs(inst[1:]...),
				opLength:  4,
			})
		case "mul":
			checkLength(inst, 4)

			tokens = append(tokens, token{
				operation: intcode.Mul,
				args:      a.getRegistersForArgs(inst[1:]...),
				opLength:  4,
			})
		case "halt":
			tokens = append(tokens, token{operation: intcode.Halt, opLength: 1,})
		default:
			fmt.Printf("unknown instruction %s, ignoring it.\n", inst[0])
		}
	}
	a.tokens = tokens
}

func (a *AsmState) Emit(autoHalt bool) (out string, err error) {
	defer func() {
		panik := recover()
		if panik != nil {
			err = errors.New(fmt.Sprint(panik))
		}
	}()
	// Preprocess and select IDs for registers
	codeLen := 0
	seenHalt := false
	var code []string
	for _, v := range a.tokens {
		codeLen += v.opLength
		if v.operation == intcode.Halt {
			seenHalt = true
		}
	}

	if !seenHalt && autoHalt {
		codeLen++
	}

	for i, reg := range a.registers { // special case for register 0
		if i == 0 || reg == nil {
			continue
		}
		reg.codeLoc = codeLen
		codeLen++
	}

	for _, reg := range a.literalRegisters {
		reg.codeLoc = codeLen
		codeLen++
	}

	for _, tok := range a.tokens {
		code = append(code, strconv.Itoa(tok.operation))
		for _, reg := range tok.args {
			if reg.id == -1 && reg.codeLoc < 0 {
				code = append(code, strconv.Itoa(codeLen+reg.codeLoc))
			} else {
				code = append(code, strconv.Itoa(reg.codeLoc))
			}
		}
	}

	if !seenHalt && autoHalt {
		code = append(code, "99")
	}

	for i, reg := range a.registers {
		if i == 0 || reg == nil {
			continue
		}
		code = append(code, strconv.Itoa(reg.content))
	}
	for _, reg := range a.literalRegisters {
		code = append(code, strconv.Itoa(reg.content))
	}
	return strings.Join(code, ","), nil
}
