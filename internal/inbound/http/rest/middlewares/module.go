package middlewares

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/inbound/http/rest/middlewares/guard"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rest",
	guard.Invoke,
)
