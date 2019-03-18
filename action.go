package devgo

import (
	"github.com/labstack/echo/v4"
)

type (
	Action func(c Context) error
)

func BeforeAction(a Action) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := GetContext(c)
			err := a(cc)
			if err != nil {
				return err
			}
			return next(cc)
		}
	}
}

func AfterAction(a Action) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := GetContext(c)
			err := next(cc)
			err2 := a(cc)
			if err != nil {
				return err
			}
			return err2
		}
	}
}

func RunAction(a Action) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := GetContext(c)
		return a(cc)
	}
}
