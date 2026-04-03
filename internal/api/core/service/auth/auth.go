package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rickferrdev/go-hasher"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
)

type Service struct {
	repository Repository
	totoken    ports.Totoken
	hasher     hasher.Hasher
	logger     *slog.Logger
}

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*domain.Consumer, error)
	Create(ctx context.Context, consumer domain.Consumer) (string, error)
}

func New(repository Repository, totoken ports.Totoken, logger *slog.Logger) *Service {
	child := logger.With(
		slog.String("location", "consumer"),
		slog.String("layer", "service"),
	)

	return &Service{
		repository: repository,
		totoken:    totoken,
		hasher:     hasher.New(hasher.DefaultCost),
		logger:     child,
	}
}

func (s *Service) Register(ctx context.Context, input ports.RegisterInput) (*ports.RegisterOutput, error) {
	child := s.logger.With(slog.String("function", "Register"))
	hash, err := s.hasher.Generate([]byte(input.Password))
	if err != nil {
		child.ErrorContext(ctx, ports.MsgAuthHashFailed,
			slog.String("email", input.Email),
			slog.Any("error", err),
		)
		return nil, ports.ErrInternalServer
	}

	id, err := s.repository.Create(ctx, domain.Consumer{
		Username: input.Username,
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: string(hash),
	})
	if err != nil {
		if errors.Is(err, ports.ErrConsumerAlreadyExists) {
			child.WarnContext(ctx, ports.MsgUserAlreadyExists, slog.String("email", input.Email))
			return nil, ports.ErrConsumerAlreadyExists
		}

		child.ErrorContext(ctx, ports.MsgSystemServiceFailed,
			slog.String("email", input.Email),
			slog.Any("error", err),
		)
		return nil, ports.ErrInternalServer
	}

	child.InfoContext(ctx, ports.MsgUserRegistered,
		slog.String("email", input.Email),
		slog.String("id", id),
	)
	return &ports.RegisterOutput{ID: id}, nil
}

func (s *Service) Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error) {
	child := s.logger.With(slog.String("function", "Login"))
	consumer, err := s.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ports.ErrConsumerNotFound) {
			child.WarnContext(ctx, ports.MsgUserNotFound, slog.String("email", input.Email))
			return nil, ports.ErrConsumerNotFound
		}

		child.ErrorContext(ctx, ports.MsgSystemServiceFailed,
			slog.String("email", input.Email),
			slog.Any("error", err),
		)
		return nil, ports.ErrInternalServer
	}

	err = s.hasher.Compare([]byte(consumer.Password), []byte(input.Password))
	if err != nil {
		child.WarnContext(ctx, ports.MsgAuthInvalidCredentials, slog.String("email", input.Email))
		return nil, ports.ErrInvalidCredentials
	}

	token, err := s.totoken.GenerateToken(consumer.ID)
	if err != nil {
		child.ErrorContext(ctx, ports.MsgAuthTokenGenFailed,
			slog.String("email", input.Email),
			slog.Any("error", err),
		)
		return nil, ports.ErrInternalServer
	}

	child.InfoContext(ctx, ports.MsgAuthLoginSuccess, slog.String("email", input.Email))
	return &ports.LoginOutput{
		Token: token,
	}, nil
}
