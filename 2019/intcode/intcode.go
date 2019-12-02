package intcode

import "fmt"

const (
	add = iota +1
	mul

	halt = 99
)

type IVM struct {
	Memory []int
	IP     int
}

func (i *IVM) Execute() {
progLoop:
	for {
		op := i.Memory[i.IP]
		switch op {
		case add:
			args := i.getArgs(3)
			ptr1, ptr2, retptr := args[0], args[1], args[2]
			i.Memory[retptr] = i.Memory[ptr1] + i.Memory[ptr2]
		case mul:
			args := i.getArgs(3)
			ptr1, ptr2, retptr := args[0], args[1], args[2]
			i.Memory[retptr] = i.Memory[ptr1] * i.Memory[ptr2]
		case halt:
			break progLoop
		default:
			panic(fmt.Sprintf("unexpected opcode %d at position %d", op, i.IP))
		}
	}
}

func (i *IVM) getArgs(count int) (args []int) {
	args = i.Memory[i.IP+1:i.IP+1+count]
	i.IP += count+1
	return
}

func (i *IVM) Print() {
	for idx, v := range i.Memory {
		if idx == i.IP {
			fmt.Printf("[%02d]", v)
		} else {
			fmt.Printf(" %02d ", v)
		}
	}
	fmt.Println()
}
