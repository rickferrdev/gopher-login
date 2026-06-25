package infra

import (
	"github.com/rickferrdev/gopher-login/internal/infra/logger"
	"github.com/rickferrdev/gopher-login/internal/infra/server"
	"github.com/rickferrdev/gopher-login/internal/infra/struct_validator"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	server.Provide,
	server.Invoke,
	struct_validator.Provide,
	logger.Provide,
)
