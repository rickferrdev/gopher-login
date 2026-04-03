package config

import (
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"go.uber.org/fx"
)

var Module = fx.Module("config", fx.Provide(
	env.New,
))
