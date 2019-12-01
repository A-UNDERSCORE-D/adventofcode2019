package main

import (
	"fmt"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

// 3325156
// part 2
// 4984866

func main() {
	inputOne := util.ReadInts("2019/01/input.txt")
	out := 0
	for _, i := range inputOne {
		out += getFuelRequired(i)
	}
	fmt.Println(out)
	fmt.Println("part 2")
	out2 := 0
	for _, i := range inputOne {
		out2 += getTotalFuelRequired(i)
	}
	fmt.Println(out2)
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
