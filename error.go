package devgo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wbsifan/devgo/errors"
)

type (
	Error interface {
		errors.Error
	}

	ErrorHandler func(err Error, c Context) error
)

var (
	catchError = defaultErrorHandler
)

func NewError(err interface{}, code ...int) Error {
	return errors.New(err, 1).SetCode(code...)
}

func IsError(err error, target error) bool {
	return errors.Is(err, target)
}

func NewErrorHandler(handler ...ErrorHandler) echo.HTTPErrorHandler {
	if len(handler) > 0 {
		catchError = handler[0]
	}
	return func(err error, c echo.Context) {
		be, ok := err.(Error)
		if !ok {
			be = errors.New(err, 0)
		}
		err = catchError(be, GetContext(c))
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func defaultErrorHandler(err Error, c Context) error {
	if Debug {
		c.Logger().Errorf("%s(%d)\n%s", err.Error(), err.Code(), err.Stack())
	} else {
		c.Logger().Errorf("%s(%d)", err.Error(), err.Code())
	}
	out := NewOutput()
	out.Code = err.Code()
	out.Message = err.Error()
	t := c.GetFormat()
	if t == FormatHtml && c.IsAjax() == false {
		return c.Display("error.html", Map{
			"path":    c.Path(),
			"code":    err.Code(),
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusBadRequest, out)
}
