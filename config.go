package devgo

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	TYPE_YAML = "YAML"
	TYPE_JSON = "JSON"
)

type (
	// AppConfig
	AppConfig struct {
		Name    string `json:"name" yaml:"name"`
		Env     string `json:"env" yaml:"env"`
		Debug   bool   `json:"debug" yaml:"debug"`
		Version string `json:"version" yaml:"version"`
	}
	// ServerConfig
	ServerConfig struct {
		Addr string `json:"addr" yaml:"addr"`
	}
	// DbConfig
	DbConfig struct {
		Driver       string `json:"driver" yaml:"driver"`
		Source       string `json:"source" yaml:"source"`
		Prefix       string `json:"prefix" yaml:"prefix"`
		MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
		MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	}
	// RedisConfig
	RedisConfig struct {
		Addr     string `json:"addr" yaml:"addr"`
		Password string `json:"password" yaml:"password"`
		Db       int    `json:"db" yaml:"db"`
	}
	JwtConfig struct {
		Name   string        `json:"name" yaml:"name"`
		Lookup string        `json:"lookup" yaml:"lookup"`
		Secret string        `json:"secret" yaml:"secret"`
		Exp    time.Duration `json:"exp" yaml:"exp"`
	}
	// Params
	Params map[string]interface{}
	// MuitlParams
	MuitlParams map[string]Params
	// ConfigParser 配置解析器接口
	ConfigParser interface {
		// Load 加载配置文件到对象
		Load(string, interface{}) error
	}

	YAMLParser struct{}
	JSONParser struct{}
)

// NewParser 创建新的解析器
func NewConfigParser(t string) ConfigParser {
	switch t {
	case TYPE_YAML:
		return &YAMLParser{}
	case TYPE_JSON:
		return &JSONParser{}
	default:
		panic("Unsupported types")
	}
}

// ParseConfigFlag 从参数获取配置文件路径
func ParseConfigFlag(name string) string {
	c := flag.String(name, "config/dev.yaml",
		"config file path \r\n default: config/dev.yaml \r\n")
	flag.Parse()
	return *c
}

//	Exists
func (this Params) Exists(key string) bool {
	_, has := this[key]
	return has
}

// Get 获取参数值
func (this Params) Get(key string) interface{} {
	v, _ := this[key]
	return v
}

// GetString 获取字符串类型
func (this Params) GetString(key string, def ...string) string {
	v, has := this[key]
	if !has {
		if len(def) > 0 {
			return def[0]
		} else {
			return ""
		}
	}
	s, _ := v.(string)
	return s
}

// GetBool
func (this Params) GetBool(key string, def ...bool) bool {
	v, has := this[key]
	if !has {
		if len(def) > 0 {
			return def[0]
		} else {
			return false
		}
	}
	s, _ := v.(bool)
	return s
}

// GetNum
func (this Params) GetInt64(key string, def ...int64) int64 {
	v, has := this[key]
	if !has {
		if len(def) > 0 {
			return def[0]
		} else {
			return 0
		}
	}
	s, _ := v.(int64)
	return s
}

// GetNum
func (this Params) GetFloat64(key string, def ...float64) float64 {
	v, has := this[key]
	if !has {
		if len(def) > 0 {
			return def[0]
		} else {
			return 0
		}
	}
	s, _ := v.(float64)
	return s
}

// GetMap
func (this Params) GetParams(key string) Params {
	v, _ := this[key].(Params)
	return v
}

// Load
func (this *YAMLParser) Load(filename string, out interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	return err
}

// Load
func (this *JSONParser) Load(filename string, out interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, out)
	return err
}
