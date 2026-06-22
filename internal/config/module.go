package config

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/config/env"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	env.Provide,
)
