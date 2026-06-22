package server

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/config/env"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)
var Invoke = fx.Invoke(Start)

func New() (*fiber.App, fiber.Router, error) {
	// Use your favorite HTTP framework; in my case, I'll write it using Fiber v3.
	app := fiber.New()

	return app, app.Group("/api/v1"), nil
}

func Start(lf fx.Lifecycle, app *fiber.App, log *slog.Logger, env *env.Env) {
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":" + env.ServerPort); err != nil {
					log.Error("error starting the HTTP server", "error", err.Error())
				}
			}()

			return nil
		},
		OnStop: app.ShutdownWithContext,
	})
}
