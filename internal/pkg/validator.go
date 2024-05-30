package pkg

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

var vtr *Validator // singleton

func GetValidator() *Validator {
	if vtr == nil {
		vtr = &Validator{v: validator.New()}
	}
	return vtr
}

// ValidateStruct validates a given struct and returns error messages if any.
func (vtr *Validator) Struct(s interface{}) (error, []string) {
	err := vtr.v.Struct(s)
	if err != nil {
		if fieldErrs, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range fieldErrs {
				errMsgs = append(errMsgs, e.Error())
			}
			return err, errMsgs
		}
		return err, []string{err.Error()}
	}

	return nil, nil
}
