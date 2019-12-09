package asm

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
)

// const fibStr = `
// ipt rCount
// add -2 rCount rCount
// add 0 1 rHigh
// add 0 1 rLow
// # loop start
// add 0 rHigh rStore
// add rLow rHigh rHigh
// add 0 rStore rLow
// add -1 rCount rCount
// jnz rCount !-4
// out rHigh
// halt
// `

// const fibStr = `
// ipt rCount
// sub rCount 2 rCount
// mov 1 rHigh
// mov 1 rLow
// add 0 rCount rCompilerLoopCounter0 # mov rCount rCompilerLoopCounter0 (for loop rCount @loop)
// add 0 . rCompilerMark@loop # mark @loop
// add 0 rHigh rStore # mov rHigh rStore
// add rLow rHigh rHigh
// add 0 rStore rLow # mov rStore rLow
// mul 1 -1 rCompilerNegStorage # sub rCompilerLoopCounter0 1 rCompilerLoopCounter0 (for loop rCount @loop)
// add rCompilerLoopCounter0 rCompilerNegStorage rCompilerLoopCounter0 # sub rCompilerLoopCounter0 1 rCompilerLoopCounter0 (for loop rCount @loop)
// jnz rCompilerLoopCounter0 rCompilerMark@loop # loop rCount @loop
// out rHigh
// halt
// `

// const testStr = `
// # register setup
// add $10 $0 %1 		# counter
// # start of loop
// add $0 %4 %2 		# mov D A
// add %2 %3 %2		# add A B A
// add $0 %3 %4 		# mov B D
// add $-1 %1 %1 	    # dec 1 c c
// jnz %1 .-16			# loop stuff
// add $0 %4 $0		# move the result to 0
// halt
// `

// const testStr = `
// mov 10 rC_loop_@loop 	# loop 10 @loop
// add 0 . rC_mark_loop 	# mov . rC_mark_loop 	# mark loop
// add 1 rCounter rCounter 	# inc rCounter
// mul -1 1 rC_negstore 	# sub rC_loop_@loop 1 rC_loop_@loop 	# loop 10 @loop
// add rC_loop_@loop rC_negstore rC_loop_@loop 	# sub rC_loop_@loop 1 rC_loop_@loop 	# loop 10 @loop
// add 0 rC_loop_@loop loop 	# jnz rC_loop_@loop loop 	# loop 10 @loop
// out rCounter
// halt
// `

const testStr = `
init rtestregister 1337
out rtestregister
halt
`

func TestTokenise(t *testing.T) {
	tokens := Tokenise(testStr)
	fmt.Println(tokens)
	assembled, regData,  err := Assemble(tokens, false) // TODO: when you have state, extract the asm so the below works
	if err != nil {
		panic(err)
	}
	fmt.Println(assembled)
	ivm := intcode.NewWithDebug(assembled, strings.Split(testStr, "\n"), regData)
	ivm.Input <- 10
	ivm.Debug = true
	ivm.HandleOut()
	fmt.Println("execute")
	ivm.Execute()
	fmt.Println("execute done")
	runtime.Gosched()
	fmt.Println(ivm.OutputData)
}
