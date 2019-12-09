package main

import (
	"fmt"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/09/input.txt")[0]
	t := time.Now()
	fmt.Println("Part 1:", part1(input), "took: ", time.Since(t))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(input), "took", time.Since(t2))
}

func runCompWithInput(prog string, input int) string {
	i := intcode.New(prog)
	i.Input <- input
	var out []int
	go i.Execute()
	for v := range i.Output {
		out = append(out, v)
	}
	return fmt.Sprint(out)
}

func part1(input string) string {
	return runCompWithInput(input, 1)
}

func part2(input string) string {
	return runCompWithInput(input, 2)
}
