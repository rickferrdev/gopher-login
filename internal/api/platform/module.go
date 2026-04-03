package platform

import (
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/api/platform/totoken"
	"github.com/rickferrdev/gopher-login/internal/api/platform/validator"
	"go.uber.org/fx"
)

var Module = fx.Module("platform", fx.Provide(
	fx.Annotate(
		totoken.New,
		fx.As(new(ports.Totoken)),
	),

	fx.Annotate(
		validator.New,
		fx.As(new(ports.Validator)),
	),
))
