package consumer

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/api/in/rest/middlewares/guard"
	"go.uber.org/fx"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

type Service interface {
	FindByUsername(ctx context.Context, username string) (*ports.ConsumerPayload, error)
	FindByID(ctx context.Context, id string) (*ports.ConsumerPayload, error)
}

type HandlerParams struct {
	fx.In
	Service Service
	Logger  *slog.Logger
	Guard   *guard.Middleware
	Router  fiber.Router
}

func New(params HandlerParams) (*Handler, error) {
	logger := params.Logger.With(
		slog.String("location", "consumer"),
		slog.String("layer", "handler"),
		slog.String("module", "consumer"),
	)
	handler := &Handler{
		service: params.Service,
		logger:  logger,
	}

	group := params.Router.Group("/consumers", params.Guard.Handle)

	group.Get("/me", handler.ObtainMe)
	group.Get("/:username", handler.ObtainConsumer)

	return handler, nil
}

func (h *Handler) ObtainMe(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	id, ok := c.Locals(guard.UserID).(string)
	if !ok || id == "" {
		return ports.NewError(
			ports.CodeAuthUnauthorized,
			ports.MessageUnauthorized,
			fiber.StatusUnauthorized,
			nil,
		)
	}

	output, err := h.service.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *Handler) ObtainConsumer(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	username := c.Params("username")
	if username == "" {
		return ports.NewError(
			ports.CodeSystemBadRequest,
			ports.MessageBadRequest,
			fiber.StatusBadRequest,
			nil,
		)
	}

	_, ok := c.Locals(guard.UserID).(string)
	if !ok {
		return ports.NewError(
			ports.CodeAuthUnauthorized,
			ports.MessageUnauthorized,
			fiber.StatusUnauthorized,
			nil,
		)
	}

	output, err := h.service.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
