package struct_validator

import (
	"github.com/go-playground/validator/v10"
	platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type StructValidator struct {
	validator validator.Validate
}

func New() (platform.StructValidator, error) {
	return &StructValidator{validator: *validator.New()}, nil
}

func (val *StructValidator) Validate(out any) error {
	return val.validator.Struct(out)
}
