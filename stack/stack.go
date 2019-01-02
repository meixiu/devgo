package stack

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
)

// stack represents a stack of program counters.
type (
	Frame uintptr
	Stack []Frame
)

const (
	Depth = 32
)

func (this Frame) String() string {
	pc := uintptr(this) - 1
	fn := runtime.FuncForPC(pc)
	s := bytes.NewBuffer([]byte{})
	if fn == nil {
		_, _ = io.WriteString(s, "unknown")
	} else {
		file, line := fn.FileLine(pc)
		fmt.Fprintf(s, "%s\n\t%s:%d", fn.Name(), file, line)
	}
	return s.String()
}

func (this Stack) String() string {
	s := make([]string, 0)
	for _, v := range this {
		s = append(s, v.String())
	}
	return strings.Join(s, "\n")
}

func New() Stack {
	var pcs [Depth]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := make(Stack, 0)
	for _, v := range pcs[0:n] {
		frames = append(frames, Frame(v))
	}
	return frames
}
