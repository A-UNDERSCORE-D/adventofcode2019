package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	start := time.Now()
	inputRange := util.GetInts(strings.Split("172930-683082", "-"))
	t := time.Now()
	c := part1(inputRange[0], inputRange[1])
	fmt.Println("Part 1:", len(c), "took:", time.Since(t))
	t = time.Now()
	fmt.Println("Part 2:", len(part2(c)), "took:", time.Since(t))
	fmt.Println("total time:", time.Since(start))
}

func part1(bottom, top int) (out []int) {
	for i := bottom; i < top; i++ {
		strNum := strconv.Itoa(i)
		if validateNum(strNum) {
			out = append(out, i)
		}
	}
	return
}

func validateNum(num string) bool {
	prev := rune(num[0])
	seenAdj := false
	for _, r := range num[1:] {
		if prev > r {
			return false
		} else if prev == r {
			seenAdj = true
		}

		prev = r
	}
	return seenAdj
}

func validatePt2(num string) bool {
	m := map[rune]int{}
	for _, r := range num {
		m[r]++
	}
	for _, v := range m {
		if v == 2 {
			return true
		}
	}
	return false
}

func part2(input []int) (out []int) {
	for _, v := range input {
		if validatePt2(strconv.Itoa(v)) {
			out = append(out, v)
		}
	}
	return
}
