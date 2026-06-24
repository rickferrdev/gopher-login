package mongo

import (
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/handlers"
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/indexes"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongo",
	handlers.Module,
	indexes.Invoke,
)
