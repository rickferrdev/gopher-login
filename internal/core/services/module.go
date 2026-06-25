package services

import (
	"github.com/rickferrdev/gopher-login/internal/core/services/user/auth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	auth.Provide,
)
