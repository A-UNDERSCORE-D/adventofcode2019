package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

const (
	black = 0
	white = 1
)

const (
	up = iota
	right
	down
	left
)

var directions = map[int]position{
	up:    {0, 1},
	right: {1, 0},
	down:  {0, -1},
	left:  {-1, 0},
}

func main() {
	t := time.Now()
	input := util.ReadLines("2019/11/input.txt")[0]
	t1 := time.Now()
	fmt.Println("Part 1:", part1(input), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(input), "took:", time.Since(t2))
	fmt.Println("total:", time.Since(t))
}

type position struct {
	x int
	y int
}

func (p *position) add(other position) position {
	return position{p.x + other.x, p.y + other.y}
}

type robit struct {
	ivm       *intcode.IVM
	pos       position
	field     map[position]int
	direction int
}

func (r *robit) turn(dir int) {
	if dir == 0 {
		dir = -1
	}
	newDir := r.direction + dir
	if newDir < up {
		newDir = left
	} else if newDir > left {
		newDir = up
	}
	r.direction = newDir
}

func newRobit(code string) *robit {
	return &robit{
		ivm:       intcode.New(code),
		pos:       position{0, 0},
		field:     make(map[position]int),
		direction: up,
	}
}

func (r *robit) run(in int) {
	r.ivm.Input <- in
	go r.ivm.Execute()
	for colour := range r.ivm.Output {
		turn := <-r.ivm.Output
		r.field[r.pos] = colour

		r.turn(turn)
		r.pos = r.pos.add(directions[r.direction])
		r.ivm.Input <- r.field[r.pos]
	}
}

func part1(input string) string {
	r := newRobit(input)
	r.run(0)
	return fmt.Sprint(len(r.field))
}

func part2(input string) string {
	r := newRobit(input)
	r.field[position{0, 0}] = 1
	r.run(1)
	var minX, minY, maxX, maxY int
	for pos, _ := range r.field {
		minX = util.Min(minX, pos.x)
		maxX = util.Max(maxX, pos.x)
		minY = util.Min(minY, pos.y)
		maxY = util.Max(maxY, pos.y)
	}

	img := image.NewGray(image.Rect(minX-10, minY-10, maxX+10, maxY+10))
	for pos, col := range r.field {
		switch col {
		case black:
			img.Set(pos.x, pos.y, color.Black)
		case white:
			img.Set(pos.x, pos.y, color.White)
		}
	}
	f, _ := os.Create("2019/11/out.png")
	defer f.Close()
	png.Encode(f, img)

	return "see picture"
}
