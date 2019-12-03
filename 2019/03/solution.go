package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

type point struct {
	x int
	y int
}

func(p *point) add(other point) point {
	return point{p.x + other.x, p.y + other.y}
}

func (p *point) distFromOrig() int {
	return util.Abs(p.x) + util.Abs(p.y)
}

type pointSlice []point
func (p pointSlice) indexOf(target point) int {
	for i, v := range p {
		if v == target {
			return i
		}
	}
	return -1
}

type set map[point]struct{}

func (s set) Union(other set) set {
	out := make(set)
	for k := range s {
		if _, exists := other[k]; exists {
			out[k] = struct{}{}
		}
	}
	return out
}

func setFromPoints(points []point) set{
	out := make(set)
	for _, p := range points {
		out[p] = struct{}{}
	}
	return out
}

var dirLookup = map[rune]point{
	'D': {-1, 0},
	'U': {1, 0},
	'L': {0, -1},
	'R': {0, 1},
}

func str2steps(steps []string) (out []point) {
	var curPos point
	out = append(out, curPos)
	for _, v := range steps {
		dirPoint := dirLookup[rune(v[0])]
		num := util.GetInt(v[1:])
		for i := 0; i < num; i++ {
			newPos := curPos.add(dirPoint)
			out = append(out, newPos)
			curPos = newPos
		}
	}
	return out
}

func main() {
	input := util.ReadLines("2019/03/input.txt")
	t := time.Now()
	var steps [][]point
	for _, i := range input {
		steps = append(steps, str2steps(strings.Split(i, ",")))
	}

	set1 := setFromPoints(steps[0])
	set2 := setFromPoints(steps[1])

	intersect := set1.Union(set2)
	minDist := math.MaxInt64
	for p := range intersect {
		if p == (point{0,0}) {
			continue
		}
		minDist = util.Min(minDist, p.distFromOrig())
	}
	fmt.Println("Part 1:", minDist, "took:", time.Since(t))

	var distances []int
	for k := range intersect {
		if k == (point{0,0}) {
			continue
		}
		distances = append(distances, pointSlice(steps[0]).indexOf(k) + pointSlice(steps[1]).indexOf(k))
	}
	fmt.Println("Part 2:", util.MinOf(distances...), "took:", time.Since(t))
}
