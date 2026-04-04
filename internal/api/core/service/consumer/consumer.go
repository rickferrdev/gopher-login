package consumer

import (
	"context"
	"log/slog"

	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
)

type Consumer struct {
	repository Repository
	logger     *slog.Logger
}

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Consumer, error)
	FindByID(ctx context.Context, id string) (*domain.Consumer, error)
}

func New(repository Repository, logger *slog.Logger) *Consumer {
	child := logger.With(
		slog.String("location", "consumer"),
		slog.String("layer", "service"),
	)

	return &Consumer{
		repository: repository,
		logger:     child,
	}
}

func (c *Consumer) FindByID(ctx context.Context, id string) (*ports.ConsumerPayload, error) {
	child := c.logger.With(slog.String("function", "FindByID"))
	if id == "" {
		child.WarnContext(ctx, ports.MsgRequestInvalidID)
		return nil, ports.ErrConsumerNotFound
	}

	consumer, err := c.repository.FindByID(ctx, id)
	if err != nil {
		if consumer == nil {
			child.WarnContext(ctx, ports.MsgUserNotFound, slog.String("id", id))
			return nil, ports.ErrConsumerNotFound
		}

		child.ErrorContext(ctx, ports.MsgSystemServiceFailed,
			slog.String("id", id),
			slog.Any("error", err),
		)
		return nil, ports.ErrInternalServer
	}

	child.DebugContext(ctx, ports.MsgUserFetchSuccess, slog.String("username", consumer.Username))
	return &ports.ConsumerPayload{
		Username: consumer.Username,
		Nickname: consumer.Nickname,
	}, nil
}

func (c *Consumer) FindByUsername(ctx context.Context, username string) (*ports.ConsumerPayload, error) {
	child := c.logger.With(slog.String("function", "FindByUsername"))
	if username == "" {
		child.WarnContext(ctx, ports.MsgRequestInvalidID)
		return nil, ports.ErrConsumerNotFound
	}

	consumer, err := c.repository.FindByUsername(ctx, username)
	if err != nil {
		child.ErrorContext(ctx, ports.MsgSystemServiceFailed,
			slog.String("username", username),
			slog.Any("error", err),
		)
		return nil, ports.ErrConsumerNotFound
	}

	if consumer == nil {
		child.WarnContext(ctx, ports.MsgUserNotFound, slog.String("username", username))
		return nil, ports.ErrConsumerNotFound
	}

	child.DebugContext(ctx, ports.MsgUserFetchSuccess, slog.String("username", consumer.Username))
	return &ports.ConsumerPayload{
		Username: consumer.Username,
		Nickname: consumer.Nickname,
	}, nil
}
