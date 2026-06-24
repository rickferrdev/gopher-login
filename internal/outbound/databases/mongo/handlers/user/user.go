package user

import (
	"context"
	"errors"

	"github.com/rickferrdev/gopher-login/internal/core/domain"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	storages "github.com/rickferrdev/gopher-login/internal/core/ports/storages"
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/constants"
	userSchema "github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/schema/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Storage struct {
	collection *mongo.Collection
}

type Params struct {
	fx.In

	Client *mongo.Client
}

func New(params Params) (storages.User, error) {
	storage := Storage{
		collection: params.Client.Database(constants.Database).Collection(constants.UserCollections),
	}

	return &storage, nil
}

func (storage *Storage) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var schema userSchema.User
	if err := storage.collection.FindOne(ctx, bson.D{{
		Key: "email", Value: email,
	}}).Decode(&schema); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ports_errors.NewDatabaseNotFound(err)
		}

		return nil, ports_errors.NewDatabaseInternal(err)
	}

	user, err := schema.ToDomain()
	if err != nil {
		return nil, ports_errors.NewDatabaseSchemaFailed(err)
	}

	return user, nil
}

func (storage *Storage) Insert(ctx context.Context, user domain.User) error {
	schema, err := userSchema.FromUserDomain(user)
	if err != nil {
		return ports_errors.NewDatabaseSchemaFailed(err)
	}

	_, err = storage.collection.InsertOne(ctx, schema)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ports_errors.NewDatabaseDuplicateKey(err)
		}

		return ports_errors.NewDatabaseInternal(err)
	}

	return nil
}
