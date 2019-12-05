package intcode

import (
	"fmt"
	"sync"
	"testing"
)

func TestIVM_Execute(t *testing.T) {
	type tst struct {
		name       string
		IVM        *IVM
		checkPtr   int
		want       int
		wantOutput bool
		input      int
		debug      bool
	}

	tsts := []tst{
		{
			name:     "test 1",
			IVM:      NewWithMem([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50}),
			checkPtr: 0,
			want:     3500,
		},
		{
			name:       "day5test",
			IVM:        NewWithMem([]int{3, 0, 4, 0, 99}),
			input:      1337,
			wantOutput: true,
			want:       1337,
		},
		{
			name:     "test sub",
			IVM:      NewWithMem([]int{1101, 100, -1, 4, 0}),
			checkPtr: 4,
			want:     99,
		},
		{
			name:       "day5 eq 8",
			IVM:        NewWithMem([]int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}),
			want:       1,
			wantOutput: true,
			input:      8,
		}, {
			name:       "day5 less than 8",
			IVM:        NewWithMem([]int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}),
			want:       1,
			wantOutput: true,
			input:      5,
		}, {
			name:       "day5 eq 8 (immed)",
			IVM:        NewWithMem([]int{3, 3, 1108, -1, 8, 3, 4, 3, 99}),
			want:       1,
			wantOutput: true,
			input:      8,
		}, {
			name:       "day5 lt 8 (immed)",
			IVM:        NewWithMem([]int{3, 3, 1107, -1, 8, 3, 4, 3, 99}),
			want:       1,
			wantOutput: true,
			input:      5,
		}, {
			name: "day5 check8",
			IVM: NewWithMem([]int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}),
			want:       999,
			wantOutput: true,
			input:      5,
		}, {
			name: "day5 check8 2",
			IVM: NewWithMem([]int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}),
			want:       1000,
			wantOutput: true,
			input:      8,
		}, {
			name: "day5 check8 3",
			IVM: NewWithMem([]int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
				1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
				999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}),
			want:       1001,
			wantOutput: true,
			input:      9,
		}, {
			name:       "nonzero input",
			IVM:        NewWithMem([]int{3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9}),
			checkPtr:   0,
			want:       0,
			wantOutput: true,
			input:      0,
		},
		{
			name:       "nonzero input (immed)",
			IVM:        NewWithMem([]int{3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9}),
			checkPtr:   0,
			want:       0,
			wantOutput: true,
			input:      0,
		},
	}
	for _, tt := range tsts {
		t.Run(tt.name, func(t *testing.T) {
			tt.IVM.Debug = tt.debug
			tt.IVM.Input <- tt.input
			var out []int
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for res := range tt.IVM.Output {
					out = append(out, res)
				}
				wg.Done()
			}()
			tt.IVM.Execute()
			wg.Wait()
			if tt.wantOutput {
				if x := out[0]; x != tt.want {
					t.Errorf("unexpected output value %d, want %d", x, tt.want)
				}
			} else {
				if res := tt.IVM.Memory[tt.checkPtr]; res != tt.want {
					t.Errorf("incorrect value at pointer %d: %d, want %d", tt.checkPtr, res, tt.want)
					tt.IVM.Print()
				}
			}
		})
	}
}

func TestGetOp(t *testing.T) {
	i := IVM{Memory: []int{1002, 4, 3, 4, 33}}
	i.Execute()
	fmt.Println(i.Memory)
}
