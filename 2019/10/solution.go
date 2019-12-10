package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

const (
	empty = iota
	asteroid
	destroyed
)

type point struct {
	x int
	y int
}

func (p *point) getSlopeFromPoint(other point) float64 {
	a, b := float64(p.y-other.y), float64(p.x-other.x)
	if a == 0 || b == 0 {
		fmt.Println(a, "/", b, "=", a/b)
	}
	return float64(p.y-other.y) / float64(p.x-other.x)
}

func (p *point) getArctan2With(other point) float64 {
	return math.Atan2(float64(p.y - other.y), float64(p.x - other.x))+1.5707963267948966 // Add to fix the rotation
}

func (p *point) String() string {
	return fmt.Sprintf("(x: %02d, y: %02d)", p.x, p.y)
}

type position struct {
	point
	content int
}

func (p *position) String() string {
	var what string
	switch p.content {
	case empty:
		what = "empty"
	case asteroid:
		what = "asteroid"
	case destroyed:
		what = "destroyed"
	}
	return fmt.Sprintf("%s at %s", what, p.point.String())
}

func main() {
	t := time.Now()
	input := util.ReadLines("2019/10/input.txt")
	t0 := time.Now()
	var positions []*position
	for y, line := range input {
		for x, p := range line {
			content := empty
			if p == '#' {
				content = asteroid
			}

			positions = append(positions, &position{
				point:   point{x, y},
				content: content,
			})
		}
	}
	fmt.Println("parse time:", time.Since(t0))
	t1 := time.Now()
	bestPos, asteroidCount := part1(positions)
	fmt.Println("Part 1:", fmt.Sprintf("%s with %d asteroids", bestPos, asteroidCount), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(positions, bestPos.point), "took:", time.Since(t2))
	fmt.Println("total:", time.Since(t))
}

type slopeSet map[float64][]int

func (s slopeSet) counterClockwiseKeyOrder() []float64 {
	switch len(s) {
	case 0:
		return nil
	case 1:
		for v := range s {
			return []float64{v}
		}
	}
	var tmp []float64
	for k := range s {
		tmp = append(tmp, k)
	}
	// sort.Sort(sort.Reverse(sort.Float64Slice(tmp)))
	// return tmp
	sort.Float64s(tmp)
	zeroIndex := -1
	// seen := false
	for i, k := range tmp {
		if k >= 0 {
			zeroIndex = i
			break
		}
	}
	if zeroIndex == -1 {
		fmt.Println("was neg1")
	}
	// return tmp
	out := make([]float64, 0, len(tmp))
	out = append(out, tmp[zeroIndex:]...)
	out = append(out, tmp[:zeroIndex]...)
	return out
}

type positionSlice []*position

func (positions positionSlice) findAtan2For(target point) slopeSet {
	out := make(slopeSet)
	for i, p := range positions {
		if p.content == empty || p.content == destroyed {
			continue
		}
		if p.point == target {
			continue
		}
		at2 := p.getArctan2With(target)
		mDist := util.Abs(p.x-target.x)+util.Abs(p.y-target.y)
		if _, exists := out[at2]; !exists {
			out[at2] = append(out[at2], i, mDist)
		} else if out[at2][1] > mDist {
			out[at2][0] = i
			out[at2][1] = mDist
		}
	}
	return out
}

func part1(positions positionSlice) (*position, int) {
	var bestPosition *position
	bestCount := math.MinInt64

	for _, pos := range positions {
		if pos.content == empty {
			continue
		}
		set := positions.findAtan2For(pos.point)
		count := len(set)
		if count > bestCount {
			bestCount = count
			bestPosition = pos
		}
	}

	return bestPosition, bestCount
}
var extreme = math.Atan2(+0.0, -1)

func part2(positions positionSlice, station point) string {
	cnt := 0
	for {
		set := positions.findAtan2For(station)
		if len(set) == 0 {
			break
		}
		for _, idx := range set.counterClockwiseKeyOrder() {
			// fmt.Println(idx, set[idx])
			cnt++
			positions[set[idx][0]].content = destroyed
			if cnt == 200 {
				return fmt.Sprintf("%02d: %s", cnt, positions[set[idx][0]])
			}
		}
	}
	return ""
}
