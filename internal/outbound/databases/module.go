package databases

import (
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"databases",
	mongo.Module,
)
