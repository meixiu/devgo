package json

import (
	stdjson "encoding/json"

	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

type (
	UnmarshalTypeError = stdjson.UnmarshalTypeError
	SyntaxError        = stdjson.SyntaxError
	RawMessage         = stdjson.RawMessage
	Number             = stdjson.Number
)

var (
	jnt           = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = jnt.Marshal
	MarshalIndent = jnt.MarshalIndent
	Unmarshal     = jnt.Unmarshal
	NewDecoder    = jnt.NewDecoder
)

func init() {
	// 开启模糊模式
	extra.RegisterFuzzyDecoders()
}
