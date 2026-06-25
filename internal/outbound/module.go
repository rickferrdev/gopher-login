package outbound

import (
	"github.com/rickferrdev/gopher-login/internal/outbound/databases"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbound",
	databases.Module,
)
