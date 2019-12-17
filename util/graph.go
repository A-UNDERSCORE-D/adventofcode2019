package util

import "fmt"

// Vec2D represents a point in space
type Vec2D struct {
	X int
	Y int
}

func (v Vec2D) Add(other Vec2D) Vec2D {
	return Vec2D{v.X + other.X, v.Y + other.Y}
}

func (v Vec2D) String() string {
	return fmt.Sprintf("(X: %d, Y: %d)", v.X, v.Y)
}

func (v Vec2D) Point() Vec2D {
	return v
}

type Node interface {
	Point() Vec2D
}

type Graph map[Vec2D]Node

func (g Graph) Filter(f func(node Node) bool) (out []Node) {
	for _, node := range g {
		if f(node) {
			out = append(out, node)
		}
	}
	return
}
