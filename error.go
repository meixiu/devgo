package devgo

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type (
	Error interface {
		Status() int
		Code() int
		Error() string
	}

	baseError struct {
		code    int
		message interface{}
	}

	ErrorHandler func(err Error, c Context) error
)

var (
	catchError = defaultErrorHandler
)

func (this *baseError) Status() int {
	he, ok := this.message.(*echo.HTTPError)
	if ok {
		return he.Code
	}
	return http.StatusOK
}

func (this *baseError) Code() int {
	return this.code
}

// Error implements of Error interface
func (this *baseError) Error() string {
	return fmt.Sprintf("%v", this.message)
}

func NewError(msg interface{}, code ...int) Error {
	be, ok := msg.(Error)
	if !ok {
		if len(code) > 0 {
			be = &baseError{
				code:    code[0],
				message: msg,
			}
		} else {
			be = &baseError{
				code:    -1,
				message: msg,
			}
		}
	}
	return be
}

func NewErrorHandler(handler ...ErrorHandler) echo.HTTPErrorHandler {
	if len(handler) > 0 {
		catchError = handler[0]
	}
	return func(err error, c echo.Context) {
		be := NewError(err)
		c.Logger().Error(be)
		err = catchError(be, GetContext(c))
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func defaultErrorHandler(err Error, c Context) error {
	out := NewOutput()
	out.Code = err.Code()
	out.Message = err.Error()
	return c.JSON(err.Status(), out)
}
