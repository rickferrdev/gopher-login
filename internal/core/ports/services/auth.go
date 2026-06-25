//go:generate mockgen -source auth.go -destination ../../../tests/mocks/ports/services.go -package ports_mocks
package ports_services

import (
	"context"
	"time"
)

type LoginOutput struct {
	Token     string
	CreatedAt time.Time
}

type RegisterOutput struct {
	UserID    string
	CreatedAt time.Time
}

type Auth interface {
	Login(ctx context.Context, email, password string) (*LoginOutput, error)
	Register(ctx context.Context, username, email, password string) (*RegisterOutput, error)
}
