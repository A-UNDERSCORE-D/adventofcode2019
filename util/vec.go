package util

type Vec2 struct {
	X int
	Y int
}

type VecSlice []Vec2

func (s VecSlice) Contains(other Vec2) bool {
	for _, v := range s {
		if v == other {
			return true
		}
	}
	return false
}

func (v *Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v *Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v *Vec2) TxDist(other Vec2) int {
	x := v.Sub(other)
	return Abs(x.X) + Abs(x.Y)
}

type Line struct {
	Start Vec2
	End   Vec2
}

// Preserves order
func getNumbersBetween(start, end int) (out []int) {
	if start < end {
		for i := start; i <= end; i++ {
			out = append(out, i)
		}
	} else {
		for i := start; i >= end; i-- {
			out = append(out, i)
		}
	}
	return
}

func (l *Line) GetStraightLinePoints() (out []Vec2) {
	if l.Start.X == l.End.X {
		for _, i := range getNumbersBetween(l.Start.Y, l.End.Y) {
			out = append(out, Vec2{X: l.Start.X, Y: i})
		}
	} else if l.Start.Y == l.End.Y {
		for _, i := range getNumbersBetween(l.Start.X, l.End.X) {
			out = append(out, Vec2{X: i, Y: l.Start.Y})
		}
	} else {
		panic("idfk") // TODO: if I want these to support angles ever, this needs to be done
	}
	return
}

func(l *Line) Intersects(other Line) (bool, Vec2) {
	pts := l.GetStraightLinePoints()
	otherpts := other.GetStraightLinePoints()
	for _, p := range pts {
		for _, op := range otherpts {
			if p == op {
				return true, p
			}
		}
	}
	return false, Vec2{}
}

type Path struct {
	Lines []Line
}

func (p Path) FindAllIntersects(other Path) []Vec2 {
	out := VecSlice{}
	for _, l := range p.Lines {
		for _, ol := range other.Lines {
			intersects, v := l.Intersects(ol)
			if !intersects {
				continue
			}
			if !out.Contains(v) {
				out = append(out, v)
			}
		}
	}
	return out
}


