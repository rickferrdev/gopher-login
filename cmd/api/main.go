package main

import (
	"github.com/rickferrdev/gopher-login/internal/config"
	"github.com/rickferrdev/gopher-login/internal/core/services"
	"github.com/rickferrdev/gopher-login/internal/inbound"
	"github.com/rickferrdev/gopher-login/internal/infra"
	"github.com/rickferrdev/gopher-login/internal/outbound"
	"github.com/rickferrdev/gopher-login/internal/platform"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		infra.Module,
		platform.Module,
		outbound.Module,
		services.Module,
		inbound.Module,
	).Run()
}
