package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Frame interface {
	File() string
	Line() int
	FullFuncName() string
	FuncName() string
}

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type pcFrame uintptr

// Pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f pcFrame) Pc() uintptr { return uintptr(f) - 1 }

// File returns the full path to the File that contains the
// function for this Frame's Pc.
func (f pcFrame) File() string {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.Pc())
	return file
}

// Line returns the Line number of source code of the
// function for this Frame's Pc.
func (f pcFrame) Line() int {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.Pc())
	return line
}

// FullFuncName returns the Name of this function, if known.
func (f pcFrame) FullFuncName() string {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source File
//	%d    source Line
//	%n    function Name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function Name and path of source File relative to the compile time
//	      GOPATH separated by \n\t (<funcname>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (f pcFrame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.FullFuncName())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.File())
		default:
			io.WriteString(s, path.Base(f.File()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.Line()))
	case 'n':
		io.WriteString(s, f.FuncName())
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText formats a stacktrace Frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f pcFrame) MarshalText() ([]byte, error) {
	name := f.FullFuncName()
	if name == "unknown" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.File(), f.Line())), nil
}

// FuncName removes the path prefix component of a function's Name reported by func.Name().
func (f pcFrame) FuncName() string {
	name := f.FullFuncName()
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
