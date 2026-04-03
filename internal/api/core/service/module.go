package service

import (
	"github.com/rickferrdev/gopher-login/internal/api/core/service/auth"
	"github.com/rickferrdev/gopher-login/internal/api/core/service/consumer"
	"go.uber.org/fx"
)

var Module = fx.Module("service", fx.Provide(
	auth.New,
	consumer.New,
))
