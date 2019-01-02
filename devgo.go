package devgo

import (
	"fmt"

	"github.com/wbsifan/devgo/json"
	yaml "gopkg.in/yaml.v2"
)

type (
	Any = interface{}
	Map = map[string]interface{}
)

// Dump
func Dump(item ...interface{}) {
	fmt.Printf("%#v", item...)
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

// JSONEncode
func JSONEncode(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}
