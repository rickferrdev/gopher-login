package env

import (
	"github.com/rickferrdev/dotenv"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Env struct {
	ServerPort string `env:"SERVER_PORT" default:"8080"`
}

func New() (*Env, error) {
	// In this usage model, github.com/rickferrdev/dotenv
	// is recommended, but for your specific use case, consider the package that best suits your needs.
	var env Env
	// or use "import "github.com/rickferrdev/dotenv/auto""
	dotenv.Collect()

	if err := dotenv.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}
