package inbound

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/inbound/http/rest"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"inbound",
	rest.Module,
)
