package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/05/input.txt")[0]
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(input string) string {
	ivm := intcode.New(input)
	ivm.Input <- 1
	var output []int
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for x := range ivm.Output {
			output = append(output, x)
		}
		wg.Done()
	}()

	ivm.Execute()
	wg.Wait()
	return strconv.Itoa(output[len(output)-1])
}

func part2(input string) string {
	ivm := intcode.New(input)
	ivm.Input <- 5
	var output []int
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for x := range ivm.Output {
			output = append(output, x)
		}
		wg.Done()
	}()

	ivm.Execute()
	wg.Wait()
	return strconv.Itoa(output[len(output)-1])
}
