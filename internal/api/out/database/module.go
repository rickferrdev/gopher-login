package database

import (
	ServiceAuth "github.com/rickferrdev/gopher-login/internal/api/core/service/auth"
	ServiceConsumer "github.com/rickferrdev/gopher-login/internal/api/core/service/consumer"
	RepositoryConsumer "github.com/rickferrdev/gopher-login/internal/api/out/database/postgres/consumer"
	"go.uber.org/fx"
)

var Module = fx.Module("database", fx.Provide(
	fx.Annotate(
		RepositoryConsumer.New,
		fx.As(new(ServiceAuth.Repository)),
		fx.As(new(ServiceConsumer.Repository)),
	),
), fx.Invoke(RepositoryConsumer.New))
