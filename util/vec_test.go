package util

import (
	"fmt"
	"testing"
)

func Test_stuff(t *testing.T) {
	v := Line{Vec2{-5, 0}, Vec2{5, 0}}
	v2 := Line{
		Start: Vec2{0, -5},
		End:   Vec2{0, 5},
	}
	fmt.Println(v.GetStraightLinePoints())
	fmt.Println(v2.GetStraightLinePoints())
	fmt.Println(v.Intersects(v2))
	fmt.Println(v2.Start.TxDist(v2.End))
}
