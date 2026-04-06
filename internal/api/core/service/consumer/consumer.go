package consumer

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"go.uber.org/fx"
)

type Consumer struct {
	repository Repository
}

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Consumer, error)
	FindByID(ctx context.Context, id string) (*domain.Consumer, error)
}

type ServiceParams struct {
	fx.In
	Repository Repository
}

func New(params ServiceParams) *Consumer {
	return &Consumer{
		repository: params.Repository,
	}
}

func (service *Consumer) FindByID(ctx context.Context, id string) (*ports.ConsumerPayload, error) {
	if id == "" {
		return nil, ports.NewError(
			ports.CodeRequestInvalidID,
			ports.MessageInvalidID,
			fiber.StatusBadRequest,
			nil,
		)
	}

	consumer, err := service.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, &ports.GopherError{Code: ports.CodeUserNotFound}) {
			return nil, ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, fiber.StatusNotFound, err)
		}

		if errors.Is(err, &ports.GopherError{Code: ports.CodeDatabaseFetchFailed}) {
			return nil, ports.NewError(ports.CodeDatabaseFetchFailed, ports.MessageInternalError, fiber.StatusInternalServerError, err)
		}

		return nil, err
	}

	return &ports.ConsumerPayload{
		Username: consumer.Username,
		Nickname: consumer.Nickname,
	}, nil
}

func (service *Consumer) FindByUsername(ctx context.Context, username string) (*ports.ConsumerPayload, error) {
	if username == "" {
		return nil, ports.NewError(
			ports.CodeSystemBadRequest,
			ports.MessageBadRequest,
			fiber.StatusBadRequest,
			nil,
		)
	}

	consumer, err := service.repository.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, &ports.GopherError{Code: ports.CodeUserNotFound}) {
			return nil, ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, fiber.StatusNotFound, err)
		}

		if errors.Is(err, &ports.GopherError{Code: ports.CodeDatabaseFetchFailed}) {
			return nil, ports.NewError(ports.CodeDatabaseFetchFailed, ports.MessageInternalError, fiber.StatusInternalServerError, err)
		}

		return nil, err
	}

	return &ports.ConsumerPayload{
		Username: consumer.Username,
		Nickname: consumer.Nickname,
	}, nil
}
