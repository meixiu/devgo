package errors

import (
	"fmt"

	"github.com/meixiu/devgo/helper"
)

type (
	Error interface {
		SetCode(code ...int) *warpError
		SetPrefix(prefix ...string) *warpError
		Code() int
		Error() string
		Stack() string
		Is(target error) bool
	}

	warpError struct {
		err    error
		frame  Frame
		code   int
		prefix string
	}
)

func New(e interface{}, skip ...int) *warpError {
	var err error
	switch e := e.(type) {
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}
	return &warpError{
		err:    err,
		frame:  Caller(helper.ArgsInt(0, skip...)),
		code:   -1,
		prefix: "",
	}
}

func Is(e error, target error) bool {
	if e == target {
		return true
	}
	if e, ok := e.(*warpError); ok {
		return Is(e.err, target)
	}
	if target, ok := target.(*warpError); ok {
		return Is(e, target.err)
	}
	return false
}

func (this *warpError) Is(target error) bool {
	return Is(this, target)
}

func (this *warpError) SetPrefix(prefix ...string) *warpError {
	if len(prefix) > 0 {
		this.prefix = prefix[0]
	}
	return this
}

func (this *warpError) SetCode(code ...int) *warpError {
	if len(code) > 0 {
		this.code = code[0]
	}
	return this
}

func (this *warpError) Code() int {
	return this.code
}

func (this *warpError) Error() string {
	msg := this.err.Error()
	if this.prefix != "" {
		msg = fmt.Sprintf("%s: %s", this.prefix, msg)
	}
	return msg
}

func (this *warpError) Stack() string {
	str := fmt.Sprintf(" * %s(%d):\n   %s", this.prefix, this.code, this.frame.String())
	if err, ok := this.err.(*warpError); ok {
		str = fmt.Sprintf("%s\n%s", str, err.Stack())
	}
	return str
}
