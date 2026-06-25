package handlers

import (
	"github.com/rickferrdev/gopher-login/internal/inbound/http/rest/handlers/auth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handlers",
	auth.Invoke,
)
