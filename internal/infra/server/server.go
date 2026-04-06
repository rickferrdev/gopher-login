package server

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"go.uber.org/fx"
)

func NewServer(validator ports.Validator, logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         "Gopher Login API v1.0",
		CaseSensitive:   true,
		StrictRouting:   true,
		StructValidator: validator,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := ports.CodeSystemInternalError
			message := ports.MessageInternalError
			status := fiber.StatusInternalServerError

			var e *ports.GopherError
			var f *fiber.Error

			switch {
			case errors.As(err, &e):
				code = e.Code
				message = e.Message
				status = e.Status
			case errors.As(err, &f):
				status = f.Code
				message = ports.Message(f.Message)
			}

			log := logger.With(
				slog.String("code", string(code)),
				slog.Int("status", status),
				slog.String("path", c.Path()),
				slog.String("method", c.Method()),
				slog.String("trace", err.Error()),
			)

			switch {
			case status >= 500:
				log.ErrorContext(c.Context(), string(message))

			case status >= 400:
				if status == 429 || status == 403 {
					log.WarnContext(c.Context(), string(message))
					break
				}
				log.InfoContext(c.Context(), string(message))

			default:
				log.InfoContext(c.Context(), string(message))
			}

			return c.Status(status).JSON(fiber.Map{
				"error": message,
				"code":  code,
			})
		},
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

func RegisterLifeCycle(life fx.Lifecycle, app *fiber.App, env *env.Environment, logger *slog.Logger) {
	life.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := env.GOPHER_SERVER_PORT
				if port == "" {
					port = "3000"
				}

				go func() {
					logger.Info("server starting", slog.String("port", port), slog.String("code", string(ports.CodeServerStart)))
					if err := app.Listen("0.0.0.0:" + port); err != nil {
						logger.Error("server failed to start", slog.Any("error", err), slog.String("code", string(ports.CodeServerFailed)))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("server shutting down", slog.String("code", string(ports.CodeServerShutdown)))
				return app.ShutdownWithContext(ctx)
			},
		},
	)
}
