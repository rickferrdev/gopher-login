package config

import (
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"github.com/rickferrdev/gopher-login/internal/config/mongo"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	env.Provide,
	mongo.Provide,
)
