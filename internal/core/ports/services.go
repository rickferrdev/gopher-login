// Ports define what the application needs from the outside world
// or what the outside world can ask from the application.
//
// Prefer small interfaces, created from real use cases.
// Do not add repositories, queues, clients, or gateways here until
// your application actually needs them.
package ports

import (
	"context"

	"github.com/rickferrdev/go-ports-and-adapters-template/internal/core/domain"
)

type UserService interface {
	ObtainUsernameByID(ctx context.Context, id string) (user *domain.User, err error)
}
