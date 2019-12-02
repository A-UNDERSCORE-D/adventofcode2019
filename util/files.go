package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

// ReadEntireFile does what it says on the tin
func ReadEntireFile(name string) string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	res, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(res)
}

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

// ReadCSV reads comma separated values from a file
func ReadCSV(name string) []string {
	data := strings.ReplaceAll(ReadEntireFile(name), "\n", "")
	return strings.Split(data, ",")
}

// GetInts returns a slice of ints created from a slice of strings
func GetInts(strings []string) (out []int) {
	for _, s := range strings {
		out = append(out, GetInt(s))
	}
	return
}
