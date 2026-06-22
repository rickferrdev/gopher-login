package user

import (
	"context"

	"github.com/rickferrdev/go-ports-and-adapters-template/internal/core/domain"
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/core/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Service struct{}

func New() (ports.UserService, error) {
	return &Service{}, nil
}

func (service *Service) ObtainUsernameByID(ctx context.Context, id string) (user *domain.User, err error) {
	return &domain.User{ID: "1", Username: "rickferrdev"}, nil
}
