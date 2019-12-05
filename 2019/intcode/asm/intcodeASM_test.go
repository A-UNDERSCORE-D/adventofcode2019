package asm

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
)

const fibStr = `
ipt rCount
add -2 rCount rCount
add 0 1 rHigh
add 0 1 rLow
# loop start
add 0 rHigh rStore
add rLow rHigh rHigh
add 0 rStore rLow
add -1 rCount rCount
jnz rCount .-16
out rHigh
halt
`

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

const testStr = `
add 0 10 rCount
add 1 rNum rNum
add -1 rCount rCount
jnz rCount .-8
add 0 rNum r0
halt
`

func TestTokenise(t *testing.T) {
	tokens := Tokenise(fibStr)
	fmt.Println(tokens)
	assembled := Assemble(tokens, false)
	fmt.Println(assembled)
	ivm := intcode.New(assembled)
	ivm.Input <- 10
	ivm.HandleOut()
	ivm.Execute()
	runtime.Gosched()
	fmt.Println(ivm.OutputData)
}
