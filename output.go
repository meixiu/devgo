package devgo

type (
	output struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func NewOutput() *output {
	return &output{}
}
