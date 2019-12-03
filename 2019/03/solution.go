package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/03/input.txt")
	t := time.Now()
	fmt.Printf("Part 1: %s (took %s)\n", part1(input), time.Since(t))
	fmt.Println("Part 2:", part2(input))
	// fmt.Println("Part 2:", part2([]string{
	// 	// "R8,U5,L5,D3",
	// 	// "U7,R6,D4,L4",
	// 	// "R75,D30,R83,U83,L12,D49,R71,U7,L72",
	// 	// "U62,R66,U55,R34,D71,R55,D58,R83",
	// 	//
	// 	// "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
	// 	// "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
	// }))
}

func getAllPointsForPath(in string) []util.Vec2 {
	out := util.VecSlice{}
	curPos := util.Vec2{}
	for _, v := range strings.Split(in, ",") {
		num := util.GetInt(v[1:])
		var toAdd util.Vec2
		switch v[0] {
		case 'U':
			toAdd = util.Vec2{X: 0, Y: num}
		case 'D':
			toAdd = util.Vec2{X: 0, Y: -num}
		case 'L':
			toAdd = util.Vec2{X: -num, Y: 0}
		case 'R':
			toAdd = util.Vec2{X: num, Y: 0}
		}

		end := curPos.Add(toAdd)
		l := util.Line{Start: curPos, End: end}
		res := l.GetStraightLinePoints()
		for _, v := range res[:len(res)-1] {
			out = append(out, v)
		}
		curPos = end
	}
	return out
}

func getIntersectsForPoints(pointSet [][]util.Vec2) []util.Vec2 {
	intersects := make(map[util.Vec2][]int)
	for i, points := range pointSet {
		for _, p := range points {
			if p.X == 0 && p.Y == 0 {
				continue
			}
			if intersects[p] == nil {
				intersects[p] = make([]int, len(pointSet))
			}
			intersects[p][i]++
		}
	}
	var is []util.Vec2
	for p, i := range intersects {
		if util.MinOf(i...) > 0 {
			is = append(is, p)
		}
	}
	return is
}

func part1(input []string) string {
	var pointSet [][]util.Vec2
	for _, s := range input {
		pointSet = append(pointSet, getAllPointsForPath(s))
	}
	is := getIntersectsForPoints(pointSet)
	var intersectDists []int
	for _, v := range is {
		intersectDists = append(intersectDists, v.TxDist(util.Vec2{}))
	}
	return fmt.Sprintf("Minimum for set: %d", util.MinOf(intersectDists...))
}

func findNumberOfStepsToPos(pos util.Vec2, slice []util.Vec2) int {
	for i, v := range slice {
		if v == pos {
			return i
		}
	}

	return -1
}

func sum(in []int) (out int) {
	for _, i := range in {
		out += i
	}
	return
}

func part2(input []string) string {
	var pointSet [][]util.Vec2
	for _, s := range input {
		pointSet = append(pointSet, getAllPointsForPath(s))
	}

	intersectionDistances := make(map[util.Vec2][]int)
	is := getIntersectsForPoints(pointSet)
	for _, v := range is {
		for _, set := range pointSet {
			s := findNumberOfStepsToPos(v, set)
			if s == -1 {
				panic("asd")
			}
			intersectionDistances[v] = append(intersectionDistances[v], s)
		}
	}

	minDist := 0
	first := true
	for _, dists := range intersectionDistances {
		if first {
			minDist = sum(dists)
			first = false
		}
		minDist = util.Min(minDist, sum(dists))
	}

	return fmt.Sprint(minDist)
}
