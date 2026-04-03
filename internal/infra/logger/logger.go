package logger

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/logbull/logbull-go/logbull"
	"github.com/rickferrdev/gopher-login/internal/config/env"
)

func New(env *env.Environment) (*slog.Logger, error) {
	fmt.Printf("env: %v\n", env)
	handler, err := logbull.NewSlogHandler(logbull.Config{
		Host:      env.GOPHER_LOGBULL_HOST,
		ProjectID: env.GOPHER_LOGBULL_PROJECT_ID,
	})

	if err != nil {
		return nil, err
	}

	handler.Flush()
	time.Sleep(2 * time.Second)

	sl := slog.New(handler)

	slog.SetDefault(sl)

	return sl, nil
}
