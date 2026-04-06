package guard

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"go.uber.org/fx"
)

const UserID = "user_id"

type Middleware struct {
	totoken ports.Totoken
	logger  *slog.Logger
}

type MiddlewareParams struct {
	fx.In
	Totoken ports.Totoken
	Logger  *slog.Logger
}

func New(params MiddlewareParams) (*Middleware, error) {
	logger := params.Logger.With(
		slog.String("location", "guard"),
		slog.String("layer", "middleware"),
	)

	return &Middleware{
		totoken: params.Totoken,
		logger:  logger,
	}, nil
}

func (m *Middleware) Handle(c fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return ports.NewError(
			ports.CodeAuthUnauthorized,
			ports.MessageUnauthorized,
			fiber.StatusUnauthorized,
			nil,
		)
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ports.NewError(
			ports.CodeAuthJwtInvalidFormat,
			ports.MessageInvalidToken,
			fiber.StatusUnauthorized,
			nil,
		)
	}

	token := parts[1]

	claims, err := m.totoken.VerifyToken(token)
	if err != nil {
		return ports.NewError(
			ports.CodeAuthJwtVerifyFailed,
			ports.MessageInvalidToken,
			fiber.StatusUnauthorized,
			err,
		)
	}

	c.Locals(UserID, claims.ID)
	return c.Next()
}
