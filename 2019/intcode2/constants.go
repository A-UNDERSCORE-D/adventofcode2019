package intcode2

// opcodes
const (
	Add           = 1
	Mul           = 2
	Input         = 3
	Output        = 4
	JumpIfTrue    = 5 // essentially JNZ
	JumpIfFalse   = 6 // essentially JEZ
	LessThan      = 7
	Equals        = 8
	RelBaseOffset = 9

	Halt = 99
)

// argument modes
const (
	ModePosition  = 0
	ModeImmediate = 1
	ModeRelative  = 2
)

type opInfo struct {
	name     string
	argCount int
}

var ops = map[int]opInfo{
	Add:           {"add", 3},
	Mul:           {"mul", 3},
	Input:         {"ipt", 1},
	Output:        {"out", 1},
	JumpIfTrue:    {"jnz", 2},
	JumpIfFalse:   {"jez", 2},
	LessThan:      {"lt ", 3},
	Equals:        {"eq ", 3},
	RelBaseOffset: {"rel", 1},
	Halt:          {"hlt", 0},
}
