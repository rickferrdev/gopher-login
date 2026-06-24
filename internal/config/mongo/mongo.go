package mongo

import (
	"context"

	"github.com/rickferrdev/gopher-login/internal/config/env"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Params struct {
	fx.In

	Env  *env.Env
	Life fx.Lifecycle
}

func New(params Params) (*mongo.Client, error) {
	sv := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(params.Env.DatabaseURI).SetServerAPIOptions(sv)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, ports_errors.NewStartupFailed(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, ports_errors.NewDatabaseFailedConnect(err)
	}

	params.Life.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		if client != nil {
			return client.Disconnect(ctx)
		}

		return nil
	}})

	return client, nil
}
