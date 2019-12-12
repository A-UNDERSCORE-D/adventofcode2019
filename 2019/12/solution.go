package main

import (
	"fmt"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/12/input.txt")
	t := time.Now()
	var moons []*moon
	for _, thing := range input {
		var x, y, z int
		if _, err := fmt.Sscanf(thing, "<x=%d, y=%d, z=%d>", &x, &y, &z); err != nil {
			panic(err)
		}

		moons = append(moons, &moon{position: vec3d{x, y, z}})
	}

	t1 := time.Now()
	fmt.Println("Part 1:", part1(moons), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(copyMoonSlice(moons)), "took:", time.Since(t2))
	fmt.Println("time:", time.Since(t))
}

func copyMoonSlice(in moonSlice) moonSlice {
	var out moonSlice
	for _, m := range in {
		out = append(out, &moon{
			position: m.position,
			velocity: m.velocity,
		})
	}
	return out
}

type vec3d struct {
	x int
	y int
	z int
}

func (v *vec3d) absSum() int {
	return util.Abs(v.x) + util.Abs(v.y) + util.Abs(v.z)
}

func (v *vec3d) String() string {
	return fmt.Sprintf("%d, %d, %d", v.x, v.y, v.z)
}

func (v *vec3d) eq(other vec3d) bool {
	return v.x == other.x && v.y == other.y && v.z == other.z
}

type moon struct {
	position vec3d
	velocity vec3d
}

func getGravChanges(a, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

func (m *moon) applyGravity(others []*moon) {
	for _, other := range others {
		m.velocity.x += getGravChanges(m.position.x, other.position.x)
		m.velocity.y += getGravChanges(m.position.y, other.position.y)
		m.velocity.z += getGravChanges(m.position.z, other.position.z)
	}
}

func (m *moon) applyVelocity() {
	m.position.x += m.velocity.x
	m.position.y += m.velocity.y
	m.position.z += m.velocity.z
}

func (m *moon) potentialEnergy() int {
	return m.position.absSum()
}

func (m *moon) kineticEnergy() int {
	return m.velocity.absSum()
}

func (m *moon) totalEnergy() int {
	return m.potentialEnergy() * m.kineticEnergy()
}

func (m *moon) String() string {
	return fmt.Sprintf("{pos=<%s>, vel=<%s>, e=%d}", m.position.String(), m.velocity.String(), m.totalEnergy())
}

func (m *moon) eq(other *moon) bool {
	return m.position.eq(other.position) && m.velocity.eq(other.velocity)
}

type moonSlice []*moon

func (m moonSlice) step() {
	for _, moon := range m {
		moon.applyGravity(m)
	}
	for _, moon := range m {
		moon.applyVelocity()
	}
}

func (m moonSlice) totalEnergy() int {
	out := 0
	for _, moon := range m {
		out += moon.totalEnergy()
	}
	return out
}

func printMoons(m moonSlice) {
	for _, moon := range m {
		fmt.Println(moon)
	}
}

func part1(moons moonSlice) string {
	for i := 1; i < 1001; i++ {
		moons.step()
	}
	return fmt.Sprintf("total energy: %d", moons.totalEnergy())
}

func (m moonSlice) eq(other moonSlice) bool {
	if len(m) != len(other) {
		return false
	}
	for i, v := range m {
		if !v.eq(other[i]) {
			return false
		}
	}
	return true
}

func (m *moon) specificEq(other *moon, pos int) bool{
	switch pos {
	case 0:
		return m.position.x == other.position.x && m.velocity.x == other.velocity.x
	case 1:
		return m.position.y == other.position.y && m.velocity.y == other.velocity.y
	case 2:
		return m.position.z == other.position.z && m.velocity.z == other.velocity.z
	}
	return false
}

func part2(moons moonSlice) string {
	start := copyMoonSlice(moons)
	moonAxisEq := make([]bool, 3)
	eqFoundForAllAxis := moonAxisEq[0] && moonAxisEq[1] && moonAxisEq[2]
	counters := []int64{0,0,0}
	for !eqFoundForAllAxis {
		moons.step()
		for count := 0; count < 3; count++ {
			if !moonAxisEq[count] {
				eq := true
				for i, moon := range moons {
					if !moon.specificEq(start[i], count) {
						eq = false
					}
				}
				if eq {
					moonAxisEq[count] = true
				} else {
					counters[count]++
				}
			}
		}
		eqFoundForAllAxis = moonAxisEq[0] && moonAxisEq[1] && moonAxisEq[2]
	}
	return fmt.Sprintf("lcm is %d", util.Lcm(counters[0]+1, counters[1]+1, counters[2]+1))
}
