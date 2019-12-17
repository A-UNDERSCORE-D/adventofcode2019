package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func digits(in int) (out []int) {
	var reversed []int
	for in > 9 {
		// Left shift with division, extract with modulo
		digit := in % 10
		in = in / 10
		reversed = append(reversed, digit)
	}
	out = append(out, in)
	for i := len(reversed) - 1; i >= 0; i-- {
		out = append(out, reversed[i])
	}

	return
}

type circularList struct {
	pattern []int
	curIdx  int
}

func (r *circularList) next() int {
	if r.curIdx >= len(r.pattern) {
		r.curIdx = 0
	}
	out := r.pattern[r.curIdx]
	r.curIdx++
	return out
}

func (r *circularList) getIdx(idx int) int {
	return r.pattern[idx%len(r.pattern)]
}

func (r *circularList) reset() {
	r.curIdx = 0
}

var (
	base      = []int{0, 1, 0, -1}
	shiftBase = []int{1, 0, -1, 0}
)

func createPattern(repCount int, toRepeat []int) (out []int) {
	out = make([]int, 0, len(toRepeat)*repCount)
	for _, n := range toRepeat {
		for i := 0; i < repCount; i++ {
			out = append(out, n)
		}
	}
	return out
}

func main() {
	input := util.ReadLines("2019/16/input.txt")[0]
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(input string) string {
	res := doFFT(util.GetInts(strings.Split(input, "")), 100)
	outStr := ""
	for i := 0; i < 8; i++ {
		outStr += strconv.Itoa(res[i])
	}
	return outStr
}

func part2(input string) string {
	realSignal := strings.Repeat(input, 10000)
	toStrip := util.GetInt(realSignal[:7])
	var out []int
	for _, chr := range realSignal[toStrip:] {
		out = append(out, int(chr-'0'))
	}

	for phase := 0; phase < 100; phase++ {
		sum := 0
		for i := len(out) - 1; i >= 0; i-- {
			sum += out[i]
			out[i] = sum % 10
		}
	}

	for _, n := range out[:8] {
		fmt.Print(strconv.Itoa(n))
	}
	fmt.Println()

	//
	//
	// fmt.Println(toStrip)
	// fmt.Println(len(realSignal))
	// res := doFFTStepsFast(util.GetInts(strings.Split(realSignal, "")), toStrip, 100)
	// // res := doFFT(util.GetInts(strings.Split(realSignal, "")), 100)
	// outStr := ""
	// for i := 0; i < 8; i++ {
	// 	outStr += strconv.Itoa(res[i])
	// }
	return "outStr"
}

func toNum(inDigits []int) int {
	if len(inDigits) == 0 {
		return 0
	}
	l := len(inDigits) - 1
	out := 0
	for _, i := range inDigits {
		out += i * int(math.Pow10(l))
		l--
	}
	return out
}

func doFFT(inDigits []int, phaseCount int) []int { // FFT, no not that one
	// fmt.Println()
	for i := 0; i < phaseCount; i++ {
		fmt.Printf("\rOn iteration: %d", i)
		n := doFFTStep(inDigits)
		inDigits = n
	}
	// fmt.Println()
	return inDigits
}

func doFFTStep(input []int) []int {
	out := make([]int, 0, len(input))

	for outIdx := 1; outIdx <= len(input); outIdx++ {
		repeater := circularList{pattern: createPattern(outIdx, base), curIdx: 1}
		// fmt.Printf("%d:  ", outIdx)
		num := 0
		for _, dig := range input {
			mul := repeater.next()
			// fmt.Printf("% d*% d  + ", dig, mul)
			switch mul {
			case 1:
				num += dig
			case -1:
				num += -dig
			case 0:
			}
			// num += dig * mul
		}
		num = util.Abs(num) % 10
		// fmt.Printf(" = %d\n", num)
		out = append(out, num)
		repeater.reset()
	}
	return out
}

func doFFTStepsFast(in []int, skip int, steps int) []int {
	if !(len(in) >= 2*skip-1) {
		panic("asd")
	}

	for phase := 0; phase < steps; phase++ {
		checkSum := 0
		for _, i := range in[skip:] {
			checkSum += i
		}
		cd := digits(checkSum)
		newDigits := make([]int, 0, skip)
		newDigits = append(newDigits, cd[len(cd)-1])
		for i := skip + 2; i < len(in)+1; i++ {
			checkSum -= in[i-2]
			cd = digits(checkSum)
			newDigits = append(newDigits, cd[len(cd)-1])
		}
		in = newDigits
	}
	fmt.Println(in)
	return in
}
