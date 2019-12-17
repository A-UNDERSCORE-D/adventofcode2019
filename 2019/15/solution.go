package main

import (
	"container/list"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode2"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

var draw = false

func main() {
	t := time.Now()
	input := util.ReadLines("2019/15/input.txt")[0]
	t1 := time.Now()
	field, node, res := part1(input)
	fmt.Println("Part 1:", res, "took:", time.Since(t1))
	if draw {
		fmt.Println(`
***********************
* AND NOW FOR PART 2! *
***********************`)
		time.Sleep(time.Second * 2)
	}
	t2 := time.Now()
	fmt.Println("Part 2:", part2(field, node), "took:", time.Since(t2))
	fmt.Println("total time:", time.Since(t))
}

const (
	tileWalkable = iota
	tileUnwalkable
)

const (
	responseWall  = 0
	responseOk    = 1
	responseFound = 2
)

const (
	north = 1
	south = 2
	west  = 3
	east  = 4
)

var invertDir = map[int]int{
	north: south,
	south: north,
	west:  east,
	east:  west,
}

var dir2Point = map[int]point{
	north: {0, 1},
	south: {0, -1},
	east:  {1, 0},
	west:  {-1, 0},
}

type node struct {
	known    bool
	walkable bool
	isTarget bool
	hasO2    bool
	point
	parent *node
}

func (n node) strRep() string {
	chr := chrUnknown
	if n.known {
		chr = chrKnown
		if n.isTarget {
			chr = chrTarget
		} else if n.hasO2 {
			chr = chrO2
		} else if !n.walkable {
			chr = chrWall
		}
	}
	return string(chr)
}

func (n node) String() string {
	return fmt.Sprintf("node at %v, walkable: %t, known: %t", n.point, n.walkable, n.known)
}

type point struct {
	X int
	Y int
}

func (p point) add(other point) point {
	return point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}
func (p point) sub(other point) point {
	return point{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p point) String() string {
	return fmt.Sprintf("(X: %d, Y: %d)", p.X, p.Y)
}

type field map[point]*node

const (
	chrKnown       = '.'
	chrWall        = '▓'
	chrUnknown     = ' '
	chrTarget      = '!'
	chrDrone       = 'D'
	chrHighlighted = 'x'
	chrO2          = '*'
)

func (f field) Draw(prefix string, dronePos point, highlighted ...point) {
	out := strings.Builder{}
	out.WriteString("\033[2J")
	out.WriteString(prefix)
	out.WriteRune('\n')
	var maxX, maxY, minX, minY = math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64
	for p := range f {
		maxX = util.Max(maxX, p.X)
		minX = util.Min(minX, p.X)
		maxY = util.Max(maxY, p.Y)
		minY = util.Min(minY, p.Y)
	}
	fmt.Fprintf(&out, "Dimensions: %dx%d\n", maxX-minX, maxY-minY)
	// do Y backwards to not invert the field while printing
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			n := f[point{x, y}]
			if n == nil {
				out.WriteRune(chrUnknown)
				continue
			}
			if n.known && n.point == dronePos {
				out.WriteRune(chrDrone)
				continue
			}
			done := false
			for _, h := range highlighted {
				if n.point == h {
					out.WriteRune(chrHighlighted)
					done = true
					break
				}
			}
			if !done {
				out.WriteString(n.strRep())
			}
		}
		out.WriteRune('\n')
	}
	fmt.Println(out.String())
}

func (f field) findCommonAncestor(a, b *node) *node {
	p := a
	var aAncestors []*node
	for p != nil {
		aAncestors = append(aAncestors, p)
		p = p.parent
	}
	var common *node
	test := b
outer:
	for test != nil {
		for _, parent := range aAncestors {
			if parent == test {
				common = parent
				break outer
			}
		}
		test = test.parent
	}
	return common
}

// THIS ASSUMES B IS A PARENT OF A. DONT FUCK THIS UP
func (f field) findNodesBetween(a, b *node) []*node {
	p := a
	var nodes []*node
	// we're the last stop, sucks if there isnt any others.
	if a.parent == nil {
		return nodes
	}
	for p != b {
		if p == nil {
			break
		}
		nodes = append(nodes, p)
		p = p.parent
	}
	nodes = append(nodes, p)
	return nodes
}

type robot struct {
	ivm   *intcode2.IVM
	pos   point
	field field
}

func (r *robot) moveDir(direction int) (moved bool, foundTarget bool) {
	if direction < 1 || direction > 4 {
		panic("invalid direction given")
	}
	r.ivm.Input <- direction
	result := <-r.ivm.Output
	switch result {
	case responseWall:
		return false, false
	case responseFound:
		foundTarget = true
		fallthrough
	case responseOk:
		r.pos = r.pos.add(dir2Point[direction])
		return true, foundTarget
	default:
		panic("unexpected response from intcode computer")
	}
}

// PeakAt moves to a point, then returns back to the current point
func (r *robot) PeakAt(direction int) (walkable bool, isTarget bool) {
	ok, found := r.moveDir(direction)
	if found {
		return true, true
	}
	if ok {
		walkable = true
		if found {
			isTarget = true
		}
		// Go back to where we were, we were just looking
		r.moveDir(invertDir[direction])
	}
	return
}

func (r *robot) findPathsTo(n *node) []point {
	ancestor := r.field.findCommonAncestor(r.field[r.pos], n)
	if ancestor == nil {
		if r.pos == (point{0, 0}) {
			ancestor = r.field[r.pos]
		} else {
			fmt.Println("Could not find ancestor for things")
			return nil
		}
	}

	pathToAncestor := r.field.findNodesBetween(r.field[r.pos], ancestor)
	pathFromAncestor := r.field.findNodesBetween(n, ancestor)
	var outPath []point
	for _, v := range pathToAncestor {
		outPath = append(outPath, v.point)
	}

	for i := len(pathFromAncestor) - 1; i >= 0; i-- {
		outPath = append(outPath, pathFromAncestor[i].point)
	}
	if draw {
		// fmt.Println(outPath)
	}
	return outPath
}

func (r *robot) moveToSimple(p point) (completed bool, isTarget bool) {
	// time.Sleep(time.Millisecond)
	// r.field.Draw(r.pos, p)
	directions := r.pos.sub(p)
	movesHoriz := util.Abs(directions.X)
	movesVerti := util.Abs(directions.Y)
	if movesHoriz > 1 || movesVerti > 1 || (movesHoriz > 0 && movesVerti > 0) {
		panic("simple move is one step")
	}

	if movesHoriz == 0 && movesVerti == 0 {
		// we're trying to move nowhere, pretend it works and keep going
		return true, false
	}

	dir := -1

	if movesHoriz == 1 {
		if directions.X > 0 {
			dir = west
		} else {
			dir = east
		}
	} else {
		if directions.Y > 0 {
			dir = south
		} else {
			dir = north
		}
	}
	return r.moveDir(dir)

}

func (r *robot) moveToNode(n *node) { // THis assumes that everywhere already exists
	if draw {
		// fmt.Printf("moving from %v to %v via graph magic\n", r.pos, n.point)
	}
	if n.point == r.pos {
		return
	}
	if n.parent == nil {
		// special case, we must be at the start
		return
	}

	moves := r.findPathsTo(n)
	for _, m := range moves {
		ok, _ := r.moveToSimple(m)
		if !ok {
			panic("ran into wall with path!")
		}

	}
}

func (r *robot) breadthFirstSearch(keepGoing bool) (int, *node) {
	start := node{
		known:    true,
		walkable: true,
		point:    r.pos,
	}
	start.parent = nil
	r.field[start.point] = &start
	queue := list.New()
	queue.PushFront(&start)
	distFromStart := -1

	next := func() *node { return (queue.Remove(queue.Front())).(*node) }
	add := func(n *node) *list.Element { return queue.PushBack(n) }
	var goalLoc *node
outer:
	for queue.Len() != 0 {
		distFromStart++
		currentNode := next()
		if !currentNode.walkable {
			fmt.Printf("skipping unwalkable point %s\n", currentNode)
			continue
		}
		if draw {
			r.field.Draw("Searching for damaged oxygen generator.....", r.pos, currentNode.point)
			time.Sleep(time.Millisecond * 20)
			// fmt.Println("checking node", currentNode)
		}
		// move to the current node, look around
		r.moveToNode(currentNode)

		for i := north; i <= east; i++ {
			// Get all the nodes that have not been already seen
			curPoint := currentNode.add(dir2Point[i])
			if p := r.field[curPoint]; p != nil && p.known {
				continue
			}
			walkable, isTarget := r.PeakAt(i)
			newNode := &node{
				known:    true,
				walkable: walkable,
				isTarget: isTarget,
				point:    curPoint,
				parent:   currentNode,
			}
			r.field[curPoint] = newNode
			if isTarget {
				goalLoc = newNode
				if !keepGoing {
					break outer
				}
			}
			if walkable {
				add(newNode)
			}
		}
	}
	if draw {
		r.field.Draw("Found damaged oxygen generator!", point{-5000000000000, -5000000000000}, point{-5000000000000, -5000000000000})
	}
	loc := goalLoc
	steps := 0
	for loc.parent != nil {
		loc = loc.parent
		steps++
	}
	return steps, goalLoc
}

func part1(input string) (field, *node, string) {
	r := robot{
		ivm:   intcode2.NewFromString(input),
		pos:   point{},
		field: make(map[point]*node),
	}
	go r.ivm.Execute()
	res, endNode := r.breadthFirstSearch(true)
	endNode.hasO2 = true
	return r.field, endNode, fmt.Sprintf("Found it: %d", res)
}

var adjacents = []point{
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},
}

func (f field) findO2WithAdjacentNonO2() []*node {
	var out []*node
	for _, v := range f {
		if v.hasO2 {
			// adjCount := 0
			for _, adj := range adjacents {
				adjNode := f[v.point.add(adj)]
				if adjNode == nil || !adjNode.known || !adjNode.walkable || adjNode.hasO2 {
					continue
				}
				out = append(out, adjNode)
			}
		}
	}
	return out
}

func part2(input field, start *node) string {
	step := 0
	for {
		posses := input.findO2WithAdjacentNonO2()
		if len(posses) == 0 {
			break
		}
		for _, p := range posses {
			p.hasO2 = true
		}
		step++
		var points []point
		for _, p := range posses {
			points = append(points, p.point)
		}
		if draw {
			input.Draw("Allowing oxygen to spread......", start.point, points...)
			time.Sleep(time.Millisecond * 50)
		}
	}

	return fmt.Sprintf("Minutes for oxygen dispersal: %d", step)
}
