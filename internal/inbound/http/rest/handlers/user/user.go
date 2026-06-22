package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/core/ports"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(New)

type Handler struct {
	// I will use Router Consumption to register the routes.
	router fiber.Router

	service ports.UserService
}

func New(router fiber.Router, service ports.UserService) (*Handler, error) {
	handler := Handler{
		router:  router,
		service: service,
	}

	handler.router.Get("/users", handler.ObtainUsernameByID)

	return &handler, nil
}

type RequestObtainUsernameDTO struct {
	ID string `json:"id"`
}

type ResponseObtainUsernameDTO struct {
	Data string `json:"data"`
}

func (handler *Handler) ObtainUsernameByID(c fiber.Ctx) error {
	var body RequestObtainUsernameDTO
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"error": "mal formatted or invalid body",
		})
	}

	user, err := handler.service.ObtainUsernameByID(c.RequestCtx(), body.ID)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(map[string]any{
			"error": "error retrieving username",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseObtainUsernameDTO{
		Data: user.Username,
	})
}
