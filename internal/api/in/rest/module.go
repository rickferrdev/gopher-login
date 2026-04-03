package rest

import (
	AuthService "github.com/rickferrdev/gopher-login/internal/api/core/service/auth"
	ConsumerService "github.com/rickferrdev/gopher-login/internal/api/core/service/consumer"
	AuthHandler "github.com/rickferrdev/gopher-login/internal/api/in/rest/handler/auth"
	ConsumerHandler "github.com/rickferrdev/gopher-login/internal/api/in/rest/handler/consumer"
	GuardMiddleware "github.com/rickferrdev/gopher-login/internal/api/in/rest/middlewares/guard"
	"go.uber.org/fx"
)

var Module = fx.Module("rest", fx.Provide(
	AuthHandler.New,
	ConsumerHandler.New,

	GuardMiddleware.New,

	fx.Annotate(
		AuthService.New,
		fx.As(new(AuthHandler.Service)),
	),

	fx.Annotate(
		ConsumerService.New,
		fx.As(new(ConsumerHandler.Service)),
	),
), fx.Invoke(AuthHandler.New, ConsumerHandler.New, GuardMiddleware.New))
