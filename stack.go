package errors

import (
	"fmt"
	"runtime"
)

type Stack interface {
	StackTrace() []Frame
	StackSource() Frame
}

// Stack represents a stack of program counters.
type pcStack []uintptr

func (stack *pcStack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		for _, pc := range *stack {
			f := pcFrame(pc)
			fmt.Fprintf(st, "\n%+v", f)
		}
		fmt.Fprintf(st, "\n")
	}
}

func (stack *pcStack) StackTrace() []Frame {
	f := make([]Frame, len(*stack))
	for i := 0; i < len(f); i++ {
		f[i] = pcFrame((*stack)[i])
	}
	return f
}

func (stack *pcStack) StackSource() Frame {

	return stack.StackTrace()[0]
}

func callers(skip int) Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st pcStack = pcs[0:n]
	return &st
}
