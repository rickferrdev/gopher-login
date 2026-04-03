package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
)

type ValidatoWrapper struct {
	validator *validator.Validate
}

func New() ports.Validator {
	return &ValidatoWrapper{
		validator: validator.New(),
	}
}

func (v *ValidatoWrapper) Validate(out any) error {
	return v.validator.Struct(out)
}
