package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	t := time.Now()
	input := util.ReadLines("2019/07/input.txt")
	t1 := time.Now()
	fmt.Println("Part 1:", part1(input), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(input), "took:", time.Since(t2))
	fmt.Println("total time:", time.Since(t))
}

func runAmplifiersWithArgs(prog string, args [5]int) int {
	amps := [5]*intcode.IVM{
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
	}
	for i, ivm := range amps {
		ivm.Input <- args[i]
		go ivm.Execute()
	}
	amps[0].Input <- 0 // Amp A gets 0 as its second input
	for i, ivm := range amps[:4] {
		out := <-ivm.Output
		amps[i+1].Input <- out
	}
	out := <-amps[4].Output
	return out
}

func permutation(xs []int) (permuts [][]int) {
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int{}, a...))
		} else {
			for i := k; i < len(xs); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	rc(xs, 0)

	return permuts
}

func randNumGen(nums []int, out chan [5]int) {
	for _, v := range permutation(nums) {
		out <- [5]int{v[0], v[1], v[2], v[3], v[4]}
	}
	close(out)
}

func part1(input []string) string {
	numChan := make(chan [5]int)
	bestChan := make(chan int)
	best := math.MinInt64

	go func() {
		for b := range bestChan {
			best = util.Max(best, b)
		}
	}()

	go randNumGen([]int{0, 1, 2, 3, 4}, numChan)
	wg := sync.WaitGroup{}
	for args := range numChan {
		wg.Add(1)
		go func(a [5]int) {
			bestChan <- runAmplifiersWithArgs(input[0], a)
			wg.Done()
		}(args)
	}
	wg.Wait()
	close(bestChan)
	return fmt.Sprint(best)
}

func runAmplifiersWithArgs2(prog string, args [5]int) int {
	amps := [5]*intcode.IVM{
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
		intcode.New(prog),
	}
	for i, ivm := range amps {
		ivm.Input <- args[i]
		go ivm.Execute()
	}

	amps[0].Input <- 0 // Amp A gets 0 as its second input
	out := -1
	wg := sync.WaitGroup{}
	for i, ivm := range amps {
		wg.Add(1)
		go func(i int, ivm *intcode.IVM) {
			defer func() {
				wg.Done()
				x := recover()
				if x != nil && fmt.Sprint(x) != "send on closed channel"{
					panic(x)
				}
			}()
			for x := range ivm.Output {
				if i == 4 {
					out = x
					amps[0].Input <- x
				} else {
					amps[i+1].Input <- x
				}
			}
		}(i, ivm)
	}
	wg.Wait()
	return out
}

func part2(input []string) string {
	numChan := make(chan [5]int)
	resChan := make(chan int)
	go randNumGen([]int{5,6,7,8,9},  numChan)
	best := math.MinInt64
	go func() {
		for res := range resChan {
			best = util.Max(best, res)
		}
	}()
	wg := sync.WaitGroup{}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for v := range numChan {
				resChan <- runAmplifiersWithArgs2(input[0], v)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(resChan)
	return fmt.Sprint(best)
}
