package auth

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"go.uber.org/fx"
)

type Handler struct {
	service   Service
	validator ports.Validator
	logger    *slog.Logger
}

type Service interface {
	Register(ctx context.Context, input ports.RegisterInput) (*ports.RegisterOutput, error)
	Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error)
}

type HandlerParams struct {
	fx.In
	Service   Service
	Validator ports.Validator
	Logger    *slog.Logger
	Router    fiber.Router
}

func New(params HandlerParams) (*Handler, error) {
	logger := params.Logger.With(
		slog.String("location", "auth"),
		slog.String("layer", "handler"),
	)

	handler := &Handler{
		service:   params.Service,
		validator: params.Validator,
		logger:    logger,
	}

	group := params.Router.Group("/auth")
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

func (handler *Handler) Login(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	var body RequestLoginDTO
	if err := c.Bind().Body(&body); err != nil {
		return ports.NewError(ports.CodeRequestBindingFailed, ports.MessageBadRequest, fiber.StatusBadRequest, err)
	}

	if err := handler.validator.Validate(body); err != nil {
		return ports.NewError(ports.CodeSystemBadRequest, ports.MessageValidationFailed, fiber.StatusBadRequest, err)
	}

	input := ports.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	}

	output, err := handler.service.Login(ctx, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (handler *Handler) Register(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), ports.Timeout)
	defer cancel()

	var body RequestRegisterDTO
	if err := c.Bind().Body(&body); err != nil {
		return ports.NewError(ports.CodeRequestBindingFailed, ports.MessageBadRequest, fiber.StatusBadRequest, err)
	}

	if err := handler.validator.Validate(body); err != nil {
		return ports.NewError(ports.CodeSystemBadRequest, ports.MessageValidationFailed, fiber.StatusBadRequest, err)
	}

	input := ports.RegisterInput{
		Username: body.Username,
		Nickname: body.Nickname,
		Email:    body.Email,
		Password: body.Password,
	}

	output, err := handler.service.Register(ctx, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}
