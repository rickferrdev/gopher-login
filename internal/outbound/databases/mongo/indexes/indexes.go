package indexes

import (
	"context"
	"time"

	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/constants"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(New)

type Params struct {
	fx.In

	Client *mongo.Client
}

func New(params Params) error {
	collection := params.Client.Database(constants.Database).Collection(constants.UserCollections)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key: "email", Value: 1,
			}},
			Options: options.Index().SetUnique(true).SetName("user_email_unique"),
		},
	})

	return err
}
