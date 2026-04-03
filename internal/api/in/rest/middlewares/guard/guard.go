package guard

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
)

const UserID = "user_id"

type Middleware struct {
	totoken ports.Totoken
	child   *slog.Logger
}

func New(totoken ports.Totoken) (*Middleware, error) {
	child := slog.With(
		slog.String("location", "guard"),
		slog.String("layer", "middleware"),
	)

	return &Middleware{
		totoken: totoken,
		child:   child,
	}, nil
}

func (m *Middleware) Handle(c fiber.Ctx) error {
	group := m.child.With(
		slog.Group("request",
			slog.String("function", "Handler"),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		),
	)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	header := c.Get("Authorization")

	if header == "" {
		group.WarnContext(ctx, ports.MsgAuthUnauthorized, slog.String("error", "missing token"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		group.WarnContext(ctx, ports.MsgAuthJwtInvalidFormat, slog.String("error", "invalid token format"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
	}

	token := parts[1]

	claims, err := m.totoken.VerifyToken(token)
	if err != nil {
		group.WarnContext(ctx, ports.MsgAuthJwtVerifyFailed, slog.String("error", err.Error()))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	c.Locals(UserID, claims.ID)
	group.DebugContext(ctx, ports.MsgAuthTokenValidated, slog.String("id", claims.ID))
	return c.Next()
}
