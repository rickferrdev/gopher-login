package env

import (
	"github.com/rickferrdev/dotenv"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Env struct {
	ServerPort   string `env:"SERVER_PORT" default:"8080"`
	DatabaseURI  string `env:"DATABASE_URI" required:"true"`
	JwtSecret    string `env:"JWT_SECRET" required:"true"`
	JwtExpiresIn string `env:"JWT_EXPIRES_IN" default:"24h"`
}

func New() (*Env, error) {
	var env Env
	dotenv.Collect()

	if err := dotenv.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}
