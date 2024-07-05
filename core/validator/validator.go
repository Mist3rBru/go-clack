package validator

import "fmt"

type Validator struct {
	Name string
}

func NewValidator(name string) *Validator {
	return &Validator{
		Name: name,
	}
}

func (v *Validator) Panic(cause string) {
	panic(fmt.Sprintf("%s: %s", v.Name, cause))
}

func (v *Validator) PanicMissingParam(param string) {
	v.Panic(fmt.Sprintf("missing %sParams.%s", v.Name, param))
}

func (v *Validator) ValidateRender(render any) {
	if render == nil {
		v.PanicMissingParam("Render")
	}
}

func (v *Validator) ValidateOptions(length int) {
	if length == 0 {
		v.PanicMissingParam("Options")
	}
}
