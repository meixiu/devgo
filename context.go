package devgo

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/wbsifan/devgo/errors"
)

type (
	Context interface {
		echo.Context
		BindValidate(interface{}) error
		GetFormat() string
		SetFormat(t string)
		IsAjax() bool
		GetBody() string
		GetData() Map
		SetData(name string, value interface{})
		Display(name string, data ...interface{}) error
		RetData(data ...interface{}) error
		RetError(err interface{}, code ...int) error
	}
	context struct {
		echo.Context
	}
)

const (
	OUTPUT_DATA_KEY = "_output_data"
	RAW_BODY_KEY    = "_raw_body"
	FORMAT_TYPE_KEY = "_format_type"
	FORMAT_RAW      = "raw"
	FORMAT_HTML     = "html"
	FORMAT_JSON     = "json"
	FORMAT_JSONP    = "jsonp"
	FORMAT_XML      = "xml"
)

// BindValidate
func (this *context) BindValidate(obj interface{}) error {
	err := this.Bind(obj)
	if err != nil {
		return err
	}
	err = this.Validate(obj)
	return err
}

// GetFormat
func (this *context) GetFormat() string {
	var t = ""
	ret := this.Get(FORMAT_TYPE_KEY)
	if ret != nil {
		t = ret.(string)
	}
	return t
}

// SetFormat
func (this *context) SetFormat(t string) {
	this.Set(FORMAT_TYPE_KEY, t)
}

// IsAjax
func (this *context) IsAjax() bool {
	h := this.Request().Header.Get("X-Requested-With")
	return h == "XMLHttpRequest"
}

// GetBody
func (this *context) GetBody() string {
	reqBody := []byte{}
	if this.Request().Body != nil {
		reqBody, _ = ioutil.ReadAll(this.Request().Body)
	}
	this.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	return string(reqBody)
}

// GetData
func (this *context) GetData() Map {
	var data Map
	ret := this.Get(OUTPUT_DATA_KEY)
	if ret == nil {
		data = make(Map)
		this.Set(OUTPUT_DATA_KEY, data)
	} else {
		data = ret.(Map)
	}
	return data
}

// SetData
func (this *context) SetData(name string, value interface{}) {
	this.GetData()[name] = value
}

// Display
func (this *context) Display(name string, data ...interface{}) error {
	if len(data) > 0 {
		if viewData, ok := data[0].(Map); ok {
			for k, v := range viewData {
				this.SetData(k, v)
			}
		} else {
			return this.Render(http.StatusOK, name, data[0])
		}
	}
	return this.Render(http.StatusOK, name, this.GetData())
}

// RetData
func (this *context) RetData(data ...interface{}) error {
	out := NewOutput()
	if len(data) > 0 {
		if viewData, ok := data[0].(Map); ok {
			for k, v := range viewData {
				this.SetData(k, v)
			}
		} else {
			out.Data = data[0]
			return this.JSON(http.StatusOK, out)
		}
	}
	out.Data = this.GetData()
	return this.JSON(http.StatusOK, out)
}

// RetError
func (this *context) RetError(err interface{}, code ...int) error {
	be := errors.New(err, 1).SetCode(code...)
	return catchError(be, this)
}

// NewContext
func NewContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(GetContext(c))
		}
	}
}

func GetContext(c echo.Context) Context {
	cc, ok := c.(Context)
	if !ok {
		cc = &context{c}
	}
	return cc
}
