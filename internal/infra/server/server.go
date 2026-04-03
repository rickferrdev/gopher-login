package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"go.uber.org/fx"
)

func NewServer(validator ports.Validator) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         "Gopher Login API v1.0",
		CaseSensitive:   true,
		StrictRouting:   true,
		StructValidator: validator,
	})

	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	}))

	return app
}

func NewRouter(app *fiber.App) fiber.Router {
	return app.Group("/api/v1")
}

func RegisterLifeCycle(life fx.Lifecycle, app *fiber.App, env *env.Environment) {
	life.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := env.GOPHER_SERVER_PORT
				if port == "" {
					port = "3000"
				}

				go func() {
					slog.Info(ports.MsgServerStart, "port", port)
					if err := app.Listen("0.0.0.0:" + port); err != nil {
						slog.Error(ports.MsgServerFailed, slog.Any("error", err)) // TODO:
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				slog.Info(ports.MsgServerShutdown)
				return app.ShutdownWithContext(ctx)
			},
		},
	)
}
