package devgo

import (
	"github.com/asaskevich/govalidator"
)

type (
	Validator struct{}
)

func (this *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

func NewValidator() *Validator {
	return &Validator{}
}
