package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)
var Invoke = fx.Invoke(Start)

func New(validator platform.StructValidator) (*fiber.App, fiber.Router, error) {
	app := fiber.New(fiber.Config{
		StructValidator: validator,
		ErrorHandler:    ErrorHandler,
	})

	app.Use(recoverer.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(limiter.New())

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

func ErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	request_id := c.RequestID()

	var port *ports_errors.Error
	if errors.As(err, &port) {
		slog.ErrorContext(
			c.RequestCtx(),
			"http request failed",
			"request_id", request_id,
			"code", string(port.Code),
			"status", port.Status,
			"message", string(port.Message),
			"error", err.Error(),
		)

		return c.Status(port.Status).JSON(ports_errors.Error{
			Message: port.Message,
			Code:    port.Code,
			Status:  port.Status,
		})
	}

	var f *fiber.Error
	if errors.As(err, &f) {
		code := ports_errors.CodeHttpInternalServer
		message := ports_errors.Message(f.Message)

		if f.Code == fiber.StatusNotFound {
			code = ports_errors.CodeDatabaseNotFound
			message = ports_errors.MessageDatabaseNotFound
		}

		slog.ErrorContext(
			c.RequestCtx(),
			"http fiber error",
			"request_id", request_id,
			"code", string(code),
			"status", f.Code,
			"message", string(message),
			"error", err.Error(),
		)

		return c.Status(f.Code).JSON(ports_errors.Error{
			Message: message,
			Code:    code,
			Status:  f.Code,
		})
	}

	slog.ErrorContext(
		c.RequestCtx(),
		"unhandled http error",
		"request_id", request_id,
		"error", err.Error(),
	)

	return c.Status(500).JSON(ports_errors.Error{
		Message: ports_errors.MessageHttpInternalServer,
		Code:    ports_errors.CodeHttpInternalServer,
		Status:  500,
	})
}
