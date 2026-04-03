package env

import (
	"github.com/rickferrdev/dotenv"
)

type Environment struct {
	GOPHER_SERVER_PORT       string `env:"GOPHER_SERVER_PORT"`
	GOPHER_SERVER_JWT_SECRET string `env:"GOPHER_SERVER_JWT_SECRET"`

	GOPHER_POSTGRES_URL      string `env:"GOPHER_POSTGRES_URL"`
	GOPHER_POSTGRES_USER     string `env:"GOPHER_POSTGRES_USER"`
	GOPHER_POSTGRES_PASSWORD string `env:"GOPHER_POSTGRES_PASSWORD"`

	GOPHER_LOGBULL_PROJECT_ID string `env:"GOPHER_LOGBULL_PROJECT_ID"`
	GOPHER_LOGBULL_HOST       string `env:"GOPHER_LOGBULL_HOST"`
}

func New() (*Environment, error) {
	var env Environment

	dotenv.Collect()

	if err := dotenv.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}
