package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/A-UNDERSCORE-D/adventofcode/util"
)

func main() {
	t := time.Now()
	input := parseImage(25, 6, util.ReadLines("2019/08/input.txt")[0])
	t1 := time.Now()
	fmt.Println("Part 1:", part1(input), "took", time.Since(t1))
	t2 := time.Now()
	fmt.Println("Part 2:", part2(input), "took", time.Since(t2))
	fmt.Println("total:", time.Since(t))
}

type img [][]string

func splitStringByidx(in string, idx int) (out []string) {
	for i := 0; i < len(in)-1; i += idx {
		if i+idx > len(in) {
			out = append(out, in[i:])
		} else {
			out = append(out, in[i:i+idx])
		}
	}
	return
}

func parseImage(width, height int, data string) [][]string {
	layers := splitStringByidx(data, width*height)
	var out [][]string

	for _, layer := range layers {
		out = append(out, splitStringByidx(layer, width))
	}

	return out
}

func part1(input img) string {
	best0Count := math.MaxInt64
	var bestCount map[rune]int
	for _, layer := range input {
		count := map[rune]int{}

		for _, line := range layer {
			for _, r := range line {
				count[r]++
			}
		}

		if best0Count > count['0'] {
			best0Count = count['0']
			bestCount = count
		}
	}

	return fmt.Sprintf("%d", bestCount['1']*bestCount['2'])
}

func renderImg(in img) [][]byte {
	var out [][]byte
	for _, v := range in[0] {
		out = append(out, []byte(v))
	}
	if out == nil {
		panic("out is nil")
	}

	for _, layer := range in[1:] {
		for j, line := range layer {
			for l, r := range line {
				outPix := out[j][l]
				switch outPix {
				case '2': // transparent
					out[j][l] = byte(r)
				}
			}
		}
	}

	return out
}

func part2(input img) string {
	img := renderImg(input)
	outImage := image.NewRGBA(image.Rect(0, 0, 25, 6))
	for y, layer := range img {
		for x, r := range layer {
			switch r {
			case '0':
				outImage.Set(x, y, color.Black)
			case '1':
				outImage.Set(x, y, color.White)
			}
		}
	}
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, outImage)
	return "stuff2"
}
