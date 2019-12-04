package main

import (
	"fmt"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

// 3325156
// part 2
// 4984866

func main() {
	start := time.Now()
	inputOne := util.ReadInts("2019/01/input.txt")
	out := 0
	t := time.Now()
	for _, i := range inputOne {
		out += getFuelRequired(i)
	}
	fmt.Println("part 1", out, "took", time.Since(t))
	t = time.Now()
	out2 := 0
	for _, i := range inputOne {
		out2 += getTotalFuelRequired(i)
	}
	fmt.Println("part 2", out2, "took", time.Since(t))
	fmt.Println("total:", time.Since(start))
}

func getTotalFuelRequired(mass int) int {
	inMassFuel := getFuelRequired(mass)
	fuelForFuel := getFuelRequired(inMassFuel)
	out := inMassFuel + fuelForFuel
	for {
		f := getFuelRequired(fuelForFuel)
		if f <= 0 {
			break
		}
		fuelForFuel = f
		out += f
	}

	return out
}

func getFuelRequired(mass int) int {
	return util.RoundDown((float64(mass) / 3) -2)
}
