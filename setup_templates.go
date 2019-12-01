package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"time"
)

const mainFile = `package main

import (
	"fmt"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("{{.Year}}/{{.Day}}/input.txt")
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(input []string) string {
	return "stuff"
}

func part2(input []string) string {
	return "stuff2"
}
`

func main() {
	var day, year int
	flag.IntVar(&day, "day", -1, "The day to create the template for")
	flag.IntVar(&year, "year", -1, "The year to create the template for")
	flag.Parse()
	if day == -1 {
		day = time.Now().Day()
		fmt.Printf("assuming day is %02d\n", day)
	}

	if year == -1 {
		year = time.Now().Year()
		fmt.Printf("assuming year is %02d\n", year)
	}

	if err := os.MkdirAll(fmt.Sprintf("./%d/%02d", year, day), 0o755); err != nil {
		panic(err)
	}

	t := template.Must(template.New("f").Parse(mainFile))

	f, err := os.OpenFile(fmt.Sprintf("./%d/%02d/solution.go", year, day), os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	t.Execute(f, struct {
		Day  int
		Year int
	}{day, year})
}
