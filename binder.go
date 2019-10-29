package devgo

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/meixiu/devgo/json"
)

type (
	Binder struct {
		DefaultBinder echo.Binder
	}
)

func (this *Binder) Bind(i interface{}, c echo.Context) (err error) {
	req := c.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	// 长度大于0且是json格式
	if req.ContentLength > 0 && strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		if err = json.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return NewError(fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
			} else if se, ok := err.(*json.SyntaxError); ok {
				return NewError(fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
			} else {
				return NewError(err.Error())
			}
		}
	} else {
		if err = this.DefaultBinder.Bind(i, c); err != nil {
			return err
		}
	}
	return
}

func NewBinder() *Binder {
	return &Binder{&echo.DefaultBinder{}}
}
