package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode2"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/13/input.txt")[0]
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

type position struct {
	x int
	y int
}

const (
	tileEmpty  = 0
	tileWall   = 1
	tileBlock  = 2
	tilePaddle = 3
	tileBall   = 4
)

var destructable = map[int]bool{
	tileEmpty:  true,
	tileWall:   false,
	tileBlock:  true,
	tilePaddle: false,
	tileBall:   false,
}

type tile struct {
	Type int
}

func part1(input string) string {
	field := map[position]tile{}
	i := intcode2.NewFromString(input)
	go i.Execute()
	for x := range i.Output {
		y := <-i.Output
		typ := <-i.Output
		field[position{x, y}] = tile{typ}
	}
	count := 0
	for _, t := range field {
		if t.Type == tileBlock {
			count++
		}
	}
	fmt.Println(count)
	return "stuff"
}

func displayGame(field map[position]tile, score int) {
	maxX, maxY := math.MinInt64, math.MinInt64
	for p := range field {
		maxX = util.Max(maxX, p.x)
		maxY = util.Max(maxY, p.y)
	}
	out := &strings.Builder{}
	out.WriteString("\033c")
	fmt.Fprintf(out, "Score: %d\n", score)
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			switch field[position{x, y}].Type {
			case tileEmpty:
				out.WriteString(" ")
			case tileWall:
				out.WriteString("#")
			case tileBlock:
				out.WriteString("+")
			case tilePaddle:
				out.WriteString("_")
			case tileBall:
				out.WriteString("O")
			}
		}
		out.WriteRune('\n')
	}
	fmt.Print(out.String())
}

func playGame(ball, paddle position) int {
	switch {
	case paddle.x > ball.x:
		return -1
	case ball.x > paddle.x:
		return 1
	default:
		return 0
	}

}

func part2(input string) string {
	field := map[position]tile{}
	i := intcode2.NewFromString(input)
	i.Memory[0] = 2 // Set it to free mode
	i.Input = make(chan int, 1)
	go i.Execute()
	var ballPos, paddlePos position
	changed := make(chan struct{})
	score := 0
	checkBlocks := false
	go func() {
		for range changed {
			i.Input <- playGame(ballPos, paddlePos)
			checkBlocks = true
		}
	}()

	for {
		if checkBlocks {
			checkBlocks = false
			blockCount := 0
			for _, t := range field {
				if t.Type == tileBlock {
					blockCount++
				}
			}

			if blockCount == 0 {
				return fmt.Sprintf("score: %d", score)
			}
		}

		x := <-i.Output
		y := <-i.Output
		typ := <-i.Output
		if x == -1 && y == 0 {
			score = typ
			continue
		}

		if x ==0 && y == 0 && typ == 0 {
			break
		}

		pos := position{x, y}
		field[pos] = tile{typ}
		switch typ {
		case tileBall:
			ballPos = pos
			changed <- struct{}{}
			displayGame(field, score)
			time.Sleep(time.Millisecond*5)
		case tilePaddle:
			paddlePos = pos
		}
	}
	displayGame(field, score)

	return fmt.Sprintf("score: %d", score)
}
