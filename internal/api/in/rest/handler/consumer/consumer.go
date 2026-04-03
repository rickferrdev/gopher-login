package consumer

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/api/in/rest/middlewares/guard"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

type Service interface {
	FindByUsername(ctx context.Context, username string) (*ports.ConsumerPayload, error)
	FindByID(ctx context.Context, id string) (*ports.ConsumerPayload, error)
}

func New(router fiber.Router, logger *slog.Logger, guard *guard.Middleware, service Service) (*Handler, error) {
	child := logger.With(
		slog.String("location", "consumer"),
		slog.String("layer", "handler"),
		slog.String("module", "consumer"),
	)
	handler := &Handler{
		service: service,
		logger:  child,
	}

	group := router.Group("/consumers", guard.Handle)

	group.Get("/me", handler.ObtainMe)
	group.Get("/:username", handler.ObtainConsumer)

	return handler, nil
}

func (h *Handler) ObtainMe(c fiber.Ctx) error {
	group := h.logger.With(
		slog.String("function", "ObtainMe"),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
	)

	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	id, ok := c.Locals(guard.UserID).(string)
	if !ok || id == "" {
		group.WarnContext(ctx, ports.MsgAuthUnauthorized, slog.String("id", id), slog.String("error", "malformatted or invalid locations"))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	output, err := h.service.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrConsumerNotFound) {
			group.WarnContext(ctx, ports.MsgUserNotFound, slog.String("id", id), slog.Any("error", err.Error()))
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "consumer not found"})
		}
		group.ErrorContext(ctx, ports.MsgSystemInternalError, slog.String("id", id), slog.String("error", err.Error()))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not process request"})
	}

	group.InfoContext(ctx, ports.MsgUserFetchSuccess, slog.String("id", id))
	return c.Status(http.StatusOK).JSON(output)
}

func (h *Handler) ObtainConsumer(c fiber.Ctx) error {
	child := h.logger.With(
		slog.String("function", "ObtainConsumer"),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
	)

	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	username := c.Params("username")
	if username == "" {
		child.WarnContext(ctx, ports.MsgSystemBadRequest, slog.String("username", username), slog.String("error", "username parameter is required"))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "username parameter is required"})
	}

	id, ok := c.Locals(guard.UserID).(string)
	if !ok || id == "" {
		child.WarnContext(ctx, ports.MsgAuthUnauthorized, slog.String("id", id), slog.String("error", "malformatted or invalid locations"))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	output, err := h.service.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ports.ErrConsumerNotFound) {
			child.WarnContext(ctx, ports.MsgUserNotFound, slog.String("username", username), slog.String("error", err.Error()))
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "consumer not found"})
		}
		child.ErrorContext(ctx, ports.MsgSystemInternalError, slog.String("username", username), slog.String("error", err.Error()))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not process request"})
	}

	child.InfoContext(ctx, ports.MsgUserFetchSuccess, slog.String("username", username))
	return c.Status(http.StatusOK).JSON(output)
}
