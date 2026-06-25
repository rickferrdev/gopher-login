package rest

import (
	"github.com/rickferrdev/gopher-login/internal/inbound/http/rest/handlers"
	"github.com/rickferrdev/gopher-login/internal/inbound/http/rest/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rest",
	handlers.Module,
	middlewares.Module,
)
