package intcode

import "testing"

func TestIVM_Execute(t *testing.T) {
	type tst struct {
		name string
		IVM
		checkPtr int
		want     int
	}

	tsts := []tst{
		{
			name: "test 1",
			IVM: IVM{
				Memory: []int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
				IP:     0,
			},
			checkPtr: 0,
			want:     3500,
		},
	}
	for _, tt := range tsts {
		t.Run(tt.name, func(t *testing.T) {
			tt.Execute()
			if res := tt.Memory[tt.checkPtr]; res != tt.want {
				t.Errorf("incorrect value at pointer %d: %d, want %d", tt.checkPtr, res, tt.want)
				tt.Print()
			}
		})
	}
}
