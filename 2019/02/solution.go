package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	fmt.Println(part1())
	fmt.Println(part2())
}

func part1() string {
	input := util.GetInts(util.ReadCSV("2019/02/input.txt"))
	t := time.Now()
	input[1] = 12
	input[2] = 2
	i := intcode.IVM{Memory:input}
	i.Execute()
	return fmt.Sprint(i.Memory[0], " took: ", time.Since(t))
}

func part2() string {
	orig := util.GetInts(util.ReadCSV("2019/02/input.txt"))

	var res []int
	t := time.Now()
	for a := 0; a < 100; a++ {
		for b := 0; b < 100; b++ {
			toRun := make([]int, len(orig))
			copy(toRun, orig)
			i := intcode.IVM{Memory: toRun}
			toRun[1] = a
			toRun[2] = b
			i.Execute()
			// out := intcode.RunIntcode(toRun)
			if i.Memory[0] == 19690720 {
				fmt.Printf("worked with: %d and %d\n", a, b)
				return fmt.Sprint(100*a + b, " took: ", time.Since(t))
			}
		}
	}

	return strconv.Itoa(res[0])
}
