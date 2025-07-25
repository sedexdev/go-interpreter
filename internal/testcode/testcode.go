package testcode

// GetProgram returns the C-- code that will be interpreted
// Code sample taken from Programming Paradigms and Languages external site
//   - Output from this code:
//   - 104 52 52 26 26 13 13 40 40 20 20 10 10 5 5 16 16 8 8 4 4 2 2 1
func GetProgram() string {
	return `
val = 104
while (val >= 2) {
	if (val % 2 == 0) {
		next = val / 2
	} else {
		next = 3 * val + 1
	}
	print val, next
	val = next
}
`
}
