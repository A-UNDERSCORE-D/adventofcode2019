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
		inputRange := util.GetInts(strings.Split("100000-999999", "-"))
		t := time.Now()
		c := part1(inputRange[0], inputRange[1])
		fmt.Println("Part 1:", len(c), "took:", time.Since(t))
		t = time.Now()
		fmt.Println("Part 2:", len(part2(c)), "took:", time.Since(t))
		fmt.Println("total time:", time.Since(start))
		t = time.Now()
		c1, c2 := part12(100000, 999999)
		fmt.Println("Ultra fast:", c1, c2, time.Since(t))
}

func part12(l, h int) (out, out2 uint16) {
	var ab, bc, cd, de, ef bool
	for a := 1; a < 10; a++ {
		for b := a; b < 10; b++ {
			for c := b; c < 10; c++ {
				for d := c; d < 10; d++ {
					for e := d; e < 10; e++ {
						for f := e; f < 10; f++ {
							num := a*100_000 + b*10_000 + c*1_000 + d*100 + e*10 + f
							if l <= num && num <= h {
								ab, bc, cd, de, ef = a == b, b == c, c == d, d == e, e == f
								if ab || bc || cd || de || ef {
									out++
								}
								if (ab && !bc) || (!ab && bc && !cd) || (!bc && cd && !de) || (!cd && de && !ef) || (!de && ef) {
									out2++
								}
							}
						}
					}
				}
			}
		}
	}
	return out, out2
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

func part1(bottom, top int) (out []int) {
	for i := bottom; i < top; i++ {
		strNum := strconv.Itoa(i)
		if validateNum(strNum) {
			out = append(out, i)
		}
	}
	return
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
