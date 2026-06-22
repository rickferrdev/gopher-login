package infra

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/infra/logger"
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/infra/server"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	server.Provide,
	server.Invoke,

	logger.Provide,
)
