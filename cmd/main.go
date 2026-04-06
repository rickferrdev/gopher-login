package main

import (
	"github.com/rickferrdev/gopher-login/internal/api/core/service"
	"github.com/rickferrdev/gopher-login/internal/api/in/rest"
	"github.com/rickferrdev/gopher-login/internal/api/out/database"
	"github.com/rickferrdev/gopher-login/internal/api/platform"
	"github.com/rickferrdev/gopher-login/internal/config"
	"github.com/rickferrdev/gopher-login/internal/infra"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		infra.Module,
		platform.Module,
		database.Module,
		service.Module,
		rest.Module,
	).Run()
}
