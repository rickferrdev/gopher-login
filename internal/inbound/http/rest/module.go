package rest

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/inbound/http/rest/handlers"
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/inbound/http/rest/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rest",
	handlers.Module,
	middlewares.Module,
)
