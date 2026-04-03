package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
)

type Handler struct {
	service   Service
	validator ports.Validator
	child     *slog.Logger
}

type Service interface {
	Register(ctx context.Context, input ports.RegisterInput) (*ports.RegisterOutput, error)
	Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error)
}

func New(group fiber.Router, service Service, validator ports.Validator, logger *slog.Logger) (*Handler, error) {
	child := logger.With(
		slog.String("location", "auth"),
		slog.String("layer", "handler"),
	)

	handler := &Handler{
		service:   service,
		validator: validator,
		child:     child,
	}

	group = group.Group("/auth")
	group.Post("/login", handler.Login)
	group.Post("/register", handler.Register)

	return handler, nil
}

type RequestLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RequestRegisterDTO struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Nickname string `json:"nickname" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *Handler) Login(c fiber.Ctx) error {
	group := h.child.With(
		slog.String("function", "Login"),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
	)

	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	var body RequestLoginDTO
	if err := c.Bind().Body(&body); err != nil {
		group.WarnContext(ctx, ports.MsgRequestBindingFailed, slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := ports.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	}

	output, err := h.service.Login(ctx, input)
	if err != nil {
		if errors.Is(err, ports.ErrConsumerNotFound) || errors.Is(err, ports.ErrInvalidCredentials) {
			group.WarnContext(ctx, ports.MsgAuthInvalidCredentials, slog.String("email", body.Email), slog.String("error", err.Error()))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
		}
		group.ErrorContext(ctx, ports.MsgSystemInternalError, slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not log in user"})
	}

	group.InfoContext(ctx, ports.MsgAuthLoginSuccess, slog.String("email", body.Email))
	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *Handler) Register(c fiber.Ctx) error {
	group := h.child.With(
		slog.String("function", "Register"),
		slog.String("path", c.Path()),
		slog.String("ip", c.IP()),
	)

	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	var body RequestRegisterDTO
	if err := c.Bind().Body(&body); err != nil {
		group.WarnContext(ctx, ports.MsgRequestBindingFailed, slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := ports.RegisterInput{
		Username: body.Username,
		Nickname: body.Nickname,
		Email:    body.Email,
		Password: body.Password,
	}

	output, err := h.service.Register(ctx, input)
	if err != nil {
		if errors.Is(err, ports.ErrConsumerAlreadyExists) {
			group.WarnContext(ctx, ports.MsgUserAlreadyExists, slog.String("email", body.Email))
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "consumer with given email or username already exists"})
		}
		if errors.Is(err, ports.ErrConsumerInvalidID) {
			group.WarnContext(ctx, ports.MsgRequestInvalidID, slog.String("email", body.Email))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid consumer ID format"})
		}

		group.ErrorContext(ctx, ports.MsgSystemServiceFailed, slog.String("email", body.Email), slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not register user"})
	}

	group.InfoContext(ctx, ports.MsgUserRegistered, slog.String("email", body.Email))
	return c.Status(fiber.StatusCreated).JSON(output)
}
