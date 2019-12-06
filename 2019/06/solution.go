package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	// f, _ := os.Create("profile")
	// pprof.StartCPUProfile(f)
	// for i := 0; i < 1000; i++ {
		input := util.ReadLines("2019/06/input.txt")
		t := time.Now()
		tree := getNetwork(input)
		fmt.Println("time to get network:", time.Since(t))
		t1 := time.Now()
		fmt.Println("Part 1:", part1(tree), "took:", time.Since(t1))
		t2 := time.Now()
		fmt.Println("Part 2:", part2(tree), "took:", time.Since(t2))
		fmt.Println("total time:", time.Since(t))
	// }
	// pprof.StopCPUProfile()
}

type orbit struct {
	name   string
	parent string
}

type tree struct {
	root    string
	forward map[string]orbit
}

func (t *tree) countDistanceTo(source, target string) int {
	if target == source {
		return 0
	}

	num := 1
	parent := t.forward[source].parent
	for {
		if parent != target {
			num++
			parent = t.forward[parent].parent
		} else {
			return num
		}
	}
}

func (t *tree) getAllAncestors(name string) []string {
	anc := t.forward[name].parent
	var out []string
	for {
		out = append(out, anc)
		if anc == t.root {
			break
		}
		anc = t.forward[anc].parent
	}
	return out
}

func (t *tree) findCommonAncestor(a, b string) orbit {
	o1, o2 := t.getAllAncestors(a), t.getAllAncestors(b)

	for _, v := range o1 {
		for _, n := range o2 {
			if v == n {
				return t.forward[n]
			}
		}
	}
	return orbit{}
}

func (t *tree) findAllWithAncestor(target string) (out []string) {
	if target == t.root {
		for k, _ := range t.forward {
			if k == t.root {
				continue
			}
			out = append(out, k)
		}
		return
	}

	for name, _ := range t.forward {
		for _, aName := range t.getAllAncestors(name) {
			if aName == target {
				out = append(out, name)
			}
		}
	}
	return
}

func (t *tree) countChildrenOfRoot() int {
	sum := 0
	for _, n := range t.forward {
		if n.name == t.root {
			continue
		}

		sum += t.countDistanceTo(n.name, t.root)
	}
	return sum
}

func (t *tree) countChildren(name string) int {
	all := t.findAllWithAncestor(name)
	sum := 0
	for _, n := range all {
		sum += t.countDistanceTo(n, name)
	}
	return sum
}

func getNetwork(input []string) *tree {
	tree := &tree{
		forward: make(map[string]orbit),
		root:    "COM",
	}

	for _, orbitStr := range input {
		idx := strings.IndexRune(orbitStr, ')')
		parent, child := orbitStr[:idx], orbitStr[idx+1:]
		tree.forward[child] = orbit{name: child, parent: parent}
	}
	return tree
}

func part1(t *tree) string {
	return fmt.Sprint(t.countChildrenOfRoot())
}

func part2(t *tree) string {
	anc := t.findCommonAncestor("YOU", "SAN")
	youDist := t.countDistanceTo("YOU", anc.name)
	sanDist := t.countDistanceTo("SAN", anc.name)
	return fmt.Sprint(youDist + sanDist - 2)
}
