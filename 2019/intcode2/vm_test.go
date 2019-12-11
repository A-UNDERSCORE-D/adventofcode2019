package intcode2

import (
	"fmt"
	"testing"

	"github.com/A-UNDERSCORE-D/adventofcode/2019/intcode"
)

func TestIVM_Execute(t *testing.T) {
	i := NewFromString("109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99")
	go i.Execute()
	for data := range i.Output {
		fmt.Printf("%d,", data)
	}
	fmt.Println()
}

const quine = "109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99"

func BenchmarkIVM_Execute(b *testing.B) {
	b.Run("new", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			i := NewFromString(quine)
			b.StartTimer()
			go i.Execute()
			for range i.Output {
			}
		}
	})
	b.Run("old", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			i := intcode.New(quine)
			b.StartTimer()
			go i.Execute()
			for range i.Output {
			}
		}
	})
}
