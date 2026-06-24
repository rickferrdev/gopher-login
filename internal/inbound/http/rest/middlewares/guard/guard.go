package guard

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gopher-login/internal/core/constants"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	ports_platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(New)

type Middleware struct {
	tokenizer ports_platform.Tokenizer
}

type Params struct {
	fx.In

	Tokenizer ports_platform.Tokenizer
}

func New(params Params) (*Middleware, error) {
	middleware := Middleware{
		tokenizer: params.Tokenizer,
	}

	return &middleware, nil
}

func (handler *Middleware) Guard(c fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if authorization == "" {
		return ports_errors.NewHttpUnauthorized(nil)
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == authorization {
		return ports_errors.NewHttpUnauthorized(nil)
	}

	claims, err := handler.tokenizer.ValidateUserToken(token)
	if err != nil {
		return ports_errors.NewHttpUnauthorized(nil)
	}

	c.Locals(constants.UserID, claims.UserID)
	return c.Next()
}
