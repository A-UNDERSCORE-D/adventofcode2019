package util

// Min does what it says on the tin, but with ints and not float64s
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// Max does what it says on the tin, but with ints and not float64s
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of the given number
func Abs(n int) int {
	if n > 0 {
		return n
	}
	return -n
}
