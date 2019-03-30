package devgo

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	yaml "gopkg.in/yaml.v2"
)

type (
	Any = interface{}
	Map = map[string]interface{}
)

var (
	Debug = false
)

func init() {
	spew.Config.Indent = "\t"
}

func Default() *echo.Echo {
	e := echo.New()
	// setting
	e.Debug = Debug
	e.HideBanner = true
	e.HTTPErrorHandler = NewErrorHandler()
	e.Validator = NewValidator()
	e.Binder = NewBinder()
	// middleware
	e.Use(NewContext())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}][REQUEST]${remote_ip} => ${method} ${status} ` +
			`${uri} (${latency_human})` + "\n",
	}))
	// store := session.NewCookieStore([]byte(Config.Session.Secret))
	// e.Use(session.NewSession("sid", store))
	//e.Renderer = NewRenderer("views", Map{
	//	"json": JSONEncode,
	//})
	e.Logger.SetHeader(`[${time_rfc3339}][${level}]`)
	return e
}

// Dump
func Dump(item ...interface{}) {
	spew.Dump(item...)
}

// DumpJSON
func DumpJSON(item ...interface{}) {
	for _, v := range item {
		b, err := json.MarshalIndent(v, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	}
}

// DumpYAML
func DumpYAML(item ...interface{}) {
	for _, v := range item {
		b, err := yaml.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	}
}
