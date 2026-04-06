package infra

import (
	"github.com/rickferrdev/gopher-login/internal/infra/logger"
	"github.com/rickferrdev/gopher-login/internal/infra/postgres"
	"github.com/rickferrdev/gopher-login/internal/infra/server"
	"go.uber.org/fx"
)

var Module = fx.Module("infra", fx.Provide(
	logger.New,
	postgres.New,
	server.NewServer,
	server.NewRouter,
), fx.Invoke(server.RegisterLifeCycle))
