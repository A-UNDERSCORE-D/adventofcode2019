package asm

import (
	"fmt"
	"testing"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
)

func TestAsmState_Emit(t *testing.T) {
	a := AsmState{}
	a.Tokenise(`
add 10 20 %0; mul %0 22 $-3;
`)
	fmt.Println(a.registers)
	fmt.Println(a.literalRegisters)
	code, _ := a.Emit(true)
	fmt.Println(code)
	i := intcode.New(code)
	i.Execute()
	fmt.Println(i.Memory[0])
	fmt.Println(i.Memory)
}
