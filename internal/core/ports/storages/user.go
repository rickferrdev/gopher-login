//go:generate mockgen -source user.go -destination ../../../tests/mocks/ports/storage.go -package ports_mocks
package ports_storages

import (
	"context"

	"github.com/rickferrdev/gopher-login/internal/core/domain"
)

type User interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Insert(ctx context.Context, user domain.User) error
}
