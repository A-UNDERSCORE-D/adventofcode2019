package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	input := util.ReadLines("2019/06/input.txt")
	// input := []string{
	// 	"COM)B",
	// 	"B)C",
	// 	"C)D",
	// 	"D)E",
	// 	"E)F",
	// 	"B)G",
	// 	"G)H",
	// 	"D)I",
	// 	"E)J",
	// 	"J)K",
	// 	"K)L",
	// 	"K)YOU",
	// 	"I)SAN",
	// }
	t := time.Now()
	network := getNetwork(input)
	fmt.Println("time to get network:", time.Since(t))
	_ = network
	t1 := time.Now()
	fmt.Println("Part 1:", part1(network), "took:", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(network), "took:", time.Since(t2))
	fmt.Println("total time:", time.Since(t))
}

type orbit struct {
	parent *orbit
	orbits []*orbit
	name   string
}

func newOrbit(above *orbit, name string) *orbit {
	o := &orbit{
		parent: above,
		name:   name,
	}
	if above != nil {
		above.orbits = append(above.orbits, o)
	}
	return o
}

func (o *orbit) getOrbitWithName(name string) *orbit {
	if name == o.name {
		return o
	}

	for _, obt := range o.orbits {
		if x := obt.getOrbitWithName(name); x != nil {
			return x
		}
	}
	return nil
}

func (o *orbit) countDistanceTo(name string) int {
	if name == o.name {
		return 0
	}
	num := 1
	obt := o.parent
	for {
		if obt.name != name {
			num++
			obt = obt.parent
		} else {
			return num
		}
	}
}
func (o *orbit) String() string {
	return o.name
}

func (o *orbit) getAllAncestors() (out []*orbit) {
	obt := o.parent
	for obt != nil {
		out = append(out, obt)
		obt = obt.parent
	}
	return
}

func (o *orbit) findCommonAncestor(other *orbit) *orbit {
	mine := o.getAllAncestors()
	theirs := other.getAllAncestors()

	for _, ma := range mine {
		for _, v := range theirs {
			if v == ma {
				return ma
			}
		}
	}
	return nil
}

func (o *orbit) countOrbiters() int {
	sum := 0
	for _, obt := range o.orbits {
		if len(obt.orbits) != 0 {
			for _, v := range obt.orbits {
				sum += v.countOrbiters()
				sum += v.countDistanceTo("COM")
			}
		}
		sum += obt.countDistanceTo("COM")
	}
	return sum
}

type orbit2 struct {
	name   string
	parent *orbit2
}

func getNetwork(input []string) *orbit {
	baseOrbits := [][2]string{}

	for _, orbitStr := range input {
		match := strings.Split(orbitStr, ")")
		// sort.Strings(match)
		orbitee, orbiter := match[0], match[1]
		baseOrbits = append(baseOrbits, [2]string{orbitee, orbiter})
	}
	com := newOrbit(nil, "COM")

	cnt := 0
	total := 0
	for len(baseOrbits) != 0 {
		toRemove := -1
		cnt++
		if cnt == 500 {
			total += cnt
			fmt.Printf("iteration %d: %d orbits remaining\n", total, len(baseOrbits))
			cnt = 0
		}

		for i, pair := range baseOrbits {
			parent, child := pair[0], pair[1]
			if o := com.getOrbitWithName(parent); o != nil {
				toRemove = i
				newOrbit(o, child)
				break
			}
		}

		if toRemove != -1 {
			baseOrbits = baseOrbits[:toRemove+copy(baseOrbits[toRemove:], baseOrbits[toRemove+1:])]
		}
	}
	return com
}

func part1(input *orbit) string {
	return fmt.Sprint(input.countOrbiters())
}

func part2(input *orbit) string {
	you := input.getOrbitWithName("YOU")
	san := input.getOrbitWithName("SAN")
	ancestor := you.findCommonAncestor(san)
	youDist := you.countDistanceTo(ancestor.name)
	sanDist := san.countDistanceTo(ancestor.name)
	fmt.Printf("youdist: %d, sandist: %d, total %d\n", youDist, sanDist, youDist+sanDist-2)

	return fmt.Sprint(youDist + sanDist - 2)
}
