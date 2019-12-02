package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode/asm"
)

func main() {
	var code string
	var autoHalt bool
	flag.StringVar(&code, "c", "-", "Code to interpret. - reads from stdin")
	flag.BoolVar(&autoHalt, "autohalt", true, "automatically emit a halt if one does not exist")
	flag.Parse()
	if code == "-" {
		c, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("error while reading from stdin: ", err)
			os.Exit(1)
		}
		code = string(c)
	}

	a := asm.AsmState{}
	a.Tokenise(code)
	ic, err := a.Emit(autoHalt)
	if err != nil {
		fmt.Println("error while emitting intcode: ", err)
	}

	fmt.Println("Intcode for ASM:", ic)
}
