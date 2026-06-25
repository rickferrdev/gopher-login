package auth

import (
	"context"
	"strings"
	"time"

	"github.com/rickferrdev/gopher-login/internal/core/domain"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	services "github.com/rickferrdev/gopher-login/internal/core/ports/services"
	storages "github.com/rickferrdev/gopher-login/internal/core/ports/storages"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Service struct {
	storage storages.User
	jwt     platform.Tokenizer
	hasher  platform.Hasher
}

type Params struct {
	fx.In

	UserStorage    storages.User
	JwtPlatform    platform.Tokenizer
	HasherPlatform platform.Hasher
}

func New(params Params) (services.Auth, error) {
	service := Service{
		storage: params.UserStorage,
		jwt:     params.JwtPlatform,
		hasher:  params.HasherPlatform,
	}

	return &service, nil
}

func (service *Service) Login(ctx context.Context, email, password string) (*services.LoginOutput, error) {
	if email == "" || password == "" {
		return nil, ports_errors.NewHttpBadRequest(nil)
	}

	email = strings.ToLower(strings.TrimSpace(email))

	user, err := service.storage.FindByEmail(ctx, email)
	if err != nil {
		if ports_errors.IsCode(err, ports_errors.CodeDatabaseNotFound) {
			return nil, ports_errors.NewInvalidCredentials(err)
		}

		return nil, ports_errors.NewHttpInternalServer(err)
	}

	if err := service.hasher.Validate([]byte(user.Password), []byte(password)); err != nil {
		return nil, ports_errors.NewInvalidCredentials(err)
	}

	token, err := service.jwt.GenerateUserToken(user.ID)
	if err != nil {
		return nil, ports_errors.NewHttpInternalServer(err)
	}

	output := services.LoginOutput{
		Token:     token,
		CreatedAt: time.Now(),
	}

	return &output, nil
}

func (service *Service) Register(ctx context.Context, username, email, password string) (*services.RegisterOutput, error) {
	if username == "" || email == "" || password == "" {
		return nil, ports_errors.NewHttpBadRequest(nil)
	}

	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)

	hash, err := service.hasher.Generate([]byte(password))
	if err != nil {
		return nil, ports_errors.NewHttpInternalServer(err)
	}

	user := domain.User{
		ID:        bson.NewObjectID().Hex(),
		Email:     email,
		Username:  username,
		Password:  string(hash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := service.storage.Insert(ctx, user); err != nil {
		if ports_errors.IsCode(err, ports_errors.CodeDatabaseDuplicateKey) {
			return nil, ports_errors.NewUserAlreadyExists(err)
		}

		return nil, ports_errors.NewHttpInternalServer(err)
	}

	output := services.RegisterOutput{
		UserID:    user.ID,
		CreatedAt: user.CreatedAt,
	}

	return &output, nil
}
