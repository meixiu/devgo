package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type (
	Frame struct {
		frames [3]uintptr
	}
)

func Caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+2, s.frames[:])
	return s
}

func (this Frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(this.frames[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func (this Frame) String() string {
	function, file, line := this.location()
	var str string
	if function != "" {
		str += fmt.Sprintf("%s\n\t", filepath.Base(function))
	}
	if file != "" {
		str += fmt.Sprintf("%s:%d", file, line)
	}
	return str
}
