package auth

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/go-hasher"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"go.uber.org/fx"
)

type Service struct {
	repository Repository
	totoken    ports.Totoken
	hasher     hasher.Hasher
}

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*domain.Consumer, error)
	Create(ctx context.Context, consumer domain.Consumer) (string, error)
}

type ServiceParams struct {
	fx.In
	Repository Repository
	Totoken    ports.Totoken
}

func New(params ServiceParams) *Service {
	return &Service{
		repository: params.Repository,
		totoken:    params.Totoken,
		hasher:     hasher.New(hasher.DefaultCost),
	}
}

func (service *Service) Register(ctx context.Context, input ports.RegisterInput) (*ports.RegisterOutput, error) {
	hash, err := service.hasher.Generate([]byte(input.Password))
	if err != nil {
		return nil, ports.NewError(ports.CodeAuthHashFailed, ports.MessageSecurityError, fiber.StatusInternalServerError, err)
	}

	id, err := service.repository.Create(ctx, domain.Consumer{
		Username: input.Username,
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: string(hash),
	})
	if err != nil {
		if errors.Is(err, &ports.GopherError{Code: ports.CodeUserAlreadyExists}) {
			return nil, ports.NewError(ports.CodeUserAlreadyExists, ports.MessageAlreadyExists, fiber.StatusConflict, err)
		}

		if errors.Is(err, &ports.GopherError{Code: ports.CodeDatabaseCreateFailed}) {
			return nil, ports.NewError(ports.CodeDatabaseCreateFailed, ports.MessageStorageError, fiber.StatusInternalServerError, err)
		}

		return nil, err
	}

	return &ports.RegisterOutput{ID: id}, nil
}

func (service *Service) Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error) {
	consumer, err := service.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, &ports.GopherError{Code: ports.CodeUserNotFound}) {
			return nil, ports.NewError(ports.CodeAuthInvalidCredentials, ports.MessageInvalidCredentials, fiber.StatusUnauthorized, err)
		}
		return nil, err
	}

	err = service.hasher.Compare([]byte(consumer.Password), []byte(input.Password))
	if err != nil {
		return nil, ports.NewError(ports.CodeAuthInvalidCredentials, ports.MessageInvalidCredentials, fiber.StatusUnauthorized, err)
	}

	token, err := service.totoken.GenerateToken(consumer.ID)
	if err != nil {
		return nil, ports.NewError(ports.CodeAuthTokenGenFailed, ports.MessageInternalError, fiber.StatusInternalServerError, err)
	}

	return &ports.LoginOutput{
		Token: token,
	}, nil
}
