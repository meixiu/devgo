package devgo

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

// Renderer manages a pongo2 TemplateSet
type (
	Renderer struct {
		BaseDir     string
		ViewContext Map
		TplSet      *pongo2.TemplateSet
	}
)

// NewRenderer creates a new instance of Renderer
func NewRenderer(baseDir string, viewContext ...Map) *Renderer {
	fInfo, err := os.Lstat(baseDir)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if fInfo.IsDir() == false {
		log.Fatalf("%s is not a directory", baseDir)
	}

	rdr := Renderer{}
	if len(viewContext) > 0 {
		rdr.ViewContext = viewContext[0]
	}
	loader, err := pongo2.NewLocalFileSystemLoader(baseDir)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	rdr.TplSet = pongo2.NewSet("TplSet-"+filepath.Base(baseDir), loader)

	return &rdr
}

// Render implements echo.Render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// get template, compile it anf store it in cache
	tpl, err := r.TplSet.FromCache(name)
	if err != nil {
		return err
	}
	// convert supplied data to pongo2.Context
	val, err := toPongoCtx(data)
	if err != nil {
		return err
	}
	// viewData
	if len(r.ViewContext) > 0 {
		for k, v := range r.ViewContext {
			val[k] = v
		}
	}
	// generate render the template
	err = tpl.ExecuteWriter(val, w)
	return err
}

// toPongoCtx converts a pongo2.Context, struct, map[string] to
// pongo2.Context
func toPongoCtx(data interface{}) (pongo2.Context, error) {
	pc, ok := data.(map[string]interface{})
	if ok {
		return pc, nil
	}
	m := pongo2.Context{}
	v := redirectValue(reflect.ValueOf(data))
	// fmt.Println(v.Type(), v.Kind())
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			m[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	} else if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			m[k.String()] = v.MapIndex(k).Interface()
		}
	} else {
		return nil, fmt.Errorf("cant convert %T to pongo2.Context", data)
	}
	return m, nil
}

func redirectValue(value reflect.Value) reflect.Value {
	for {
		if !value.IsValid() || value.Kind() != reflect.Ptr {
			return value
		}
		res := reflect.Indirect(value)
		if res.Kind() == reflect.Ptr && value.Pointer() == res.Pointer() {
			return value
		}
		value = res
	}
}
