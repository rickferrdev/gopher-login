package handlers

import (
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/handlers/user"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongo-handlers",
	user.Provide,
)
