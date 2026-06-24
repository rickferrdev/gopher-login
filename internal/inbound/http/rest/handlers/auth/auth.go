package auth

import (
	"github.com/gofiber/fiber/v3"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	services "github.com/rickferrdev/gopher-login/internal/core/ports/services"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(New)

type Handler struct {
	router  fiber.Router
	service services.Auth
}

type Params struct {
	fx.In

	Router      fiber.Router
	AuthService services.Auth
}

func New(params Params) (*Handler, error) {
	handler := Handler{
		router:  params.Router,
		service: params.AuthService,
	}

	handler.router.Post("/auth/login", handler.Login)
	handler.router.Post("/auth/register", handler.Register)

	return &handler, nil
}

func (handler *Handler) Register(c fiber.Ctx) error {
	var body RequestRegisterDTO
	if err := c.Bind().JSON(&body); err != nil {
		return ports_errors.NewHttpBadRequest(err)
	}

	output, err := handler.service.Register(c.RequestCtx(), body.Username, body.Email, body.Password)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseRegisterDTO{
		ID:        output.UserID,
		CreatedAt: output.CreatedAt,
	})
}

func (handler *Handler) Login(c fiber.Ctx) error {
	var body RequestLoginDTO
	if err := c.Bind().JSON(&body); err != nil {
		return ports_errors.NewHttpBadRequest(err)
	}

	output, err := handler.service.Login(c.RequestCtx(), body.Email, body.Password)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ResponseLoginDTO{
		Token:     output.Token,
		CreatedAt: output.CreatedAt,
	})
}
