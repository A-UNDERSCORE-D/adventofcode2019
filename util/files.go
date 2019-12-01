package util

import (
	"bufio"
	"os"
)

// ReadLines reads all lines from the given file, or panics if the file doesnt exist
func ReadLines(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var out []string
	for s.Scan() {
		out = append(out, s.Text())
	}

	return out
}

// ReadInts is like ReadLines but returns ints instead
func ReadInts(name string) (out []int) {
	for _, l := range ReadLines(name) {
		out = append(out, GetInt(l))
	}
	return
}
