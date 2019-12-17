package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode2"
	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

var prnt = false

func main() {
	t := time.Now()
	input := util.ReadLines("2019/17/input.txt")[0]
	t1 := time.Now()
	ascii := intcode2Ascii(input)
	// findPath(ascii)
	fmt.Println("Part 1:", part1(ascii), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(ascii, input), "took:", time.Since(t2))
	fmt.Println("total:", time.Since(t))
}

const (
	robitUp = iota
	robitDown
	robitLeft
	robitRight
	scaffold
	openSpace
)

type node struct {
	util.Vec2D
	Type int
}

func (n *node) String() string {
	return fmt.Sprintf("Node at %s, type %d", n.Vec2D, n.Type)
}

func string2Graph(in string) util.Graph {
	out := make(util.Graph)
	for y, line := range strings.Split(in, "\n") {
		for x, r := range line {
			p := util.Vec2D{
				X: x,
				Y: y,
			}
			typ := -1
			switch r {
			case '.':
				typ = openSpace
			case '#':
				typ = scaffold
			case '^':
				typ = robitUp
			case 'v':
				typ = robitDown
			case '>':
				typ = robitRight
			case '<':
				typ = robitLeft
			case 'X':
				panic("whoops")
			}

			out[p] = &node{p, typ}
		}
	}
	return out
}

var dirs = map[int]util.Vec2D{
	up:    {X: 0, Y: 1},
	down:  {X: 0, Y: -1},
	left:  {X: -1, Y: 0},
	right: {X: 1, Y: 0},
}

func findIntersects(g util.Graph) []*node {
	res := g.Filter(func(n util.Node) bool {
		if n.(*node).Type != scaffold {
			return false
		}
		for _, d := range dirs {
			other := g[n.Point().Add(d)]
			if other == nil {
				return false
			}
			realNode := other.(*node)
			if realNode.Type != scaffold {
				return false
			}
		}
		return true
	})

	var out []*node
	for _, v := range res {
		out = append(out, v.(*node))
	}
	return out
}

func intcode2Ascii(input string) string {
	i := intcode2.NewFromString(input)
	out := &strings.Builder{}
	go i.Execute()
	for res := range i.Output {
		out.WriteByte(byte(res))
	}
	return out.String()
}

func part1(input string) string {
	intersects := findIntersects(string2Graph(input))
	sum := 0
	for _, i := range intersects {
		sum += i.X * i.Y
	}
	return fmt.Sprint(sum)
}

const (
	moveOdd  = "5" // A
	moveEven = "2" // B
	turn     = "R" // C

)

func draw(g util.Graph) {
	maxX, maxY, minX, minY := math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64
	for p := range g {
		maxX = util.Max(maxX, p.X)
		maxY = util.Max(maxY, p.Y)
		minX = util.Min(minX, p.X)
		minY = util.Min(minY, p.Y)
	}
	out := strings.Builder{}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			n := g[util.Vec2D{X: x, Y: y}].(*node)
			switch n.Type {
			case robitUp:
				out.WriteRune('^')
			case robitDown:
				out.WriteRune('v')
			case robitLeft:
				out.WriteRune('<')
			case robitRight:
				out.WriteRune('>')
			case scaffold:
				out.WriteRune('#')
			case openSpace:
				out.WriteRune('.')
			}
		}
		out.WriteRune('\n')
	}
	fmt.Println(out.String())
}

func readLine(ivm *intcode2.IVM) string {
	out := strings.Builder{}
	for {
		chr := <-ivm.Output
		if chr == '\n' {
			break
		}
		out.WriteRune(rune(chr))
	}
	return out.String()
}

func writeString(ivm *intcode2.IVM, str string) {
	if prnt {
		fmt.Print(readLine(ivm) + " ")
		fmt.Println("inputting: ", str)
	} else {
		readLine(ivm)
	}
	for _, r := range str {
		ivm.Input <- int(r)
	}
	ivm.Input <- int('\n')
}

const (
	up = iota
	down
	left
	right
)

// func findPath(in string) {
// 	g := string2Graph(in)
// 	robitNode := g.Filter(func(n util.Node) bool { nde := n.(*node).Type; return nde >= robitUp && nde <= robitRight })[0].(*node)
// 	running := true
// 	fmt.Println(robitNode)
// 	out := strings.Builder{}
// 	curNode := robitNode
// 	curLen := 0
// 	curDir := -1
// 	switch robitNode.Type {
// 	case robitUp:
// 		curDir = up
// 	case robitDown:
// 		curDir = down
// 	case robitLeft:
// 		curDir = left
// 	case robitRight:
// 		curDir = right
// 	}
//
// 	for running {
// 		var possible = []*node{}
// 		for dir, toAdd := range dirs {
// 			possible, ok := g[curNode.Add(toAdd)].(*node)
// 			if !ok || possible.Type != scaffold {
// 				// not usable, ignore it
// 				continue
// 			}
//
// 		}
// 	}
// }

func part2(ascii, intcode string) string {
	findPath(ascii)
	ivm := intcode2.NewFromString(intcode)
	ivm.Memory[0] = 2
	ivm.Input = make(chan int)
	go ivm.Execute()
	for i := 0; i < len(ascii); i++ {
		<-ivm.Output // skip the map
	}
	// Worked out these values by hand, sorry anyone who wanted a fancy algo
	writeString(ivm, "A,B,A,B,C,C,B,C,B,A")
	writeString(ivm, "R,12,L,8,R,12")
	writeString(ivm, "R,8,R,6,R,6,R,8")
	writeString(ivm, "R,8,L,8,R,8,R,4,R,4")
	writeString(ivm, "n")
	for i := 0; i < len(ascii)+1; i++ {
		if prnt {
			// skip the map
			fmt.Print(string(<-ivm.Output))
		} else {
			<-ivm.Output
		}
	}
	return fmt.Sprint(<-ivm.Output)
}
