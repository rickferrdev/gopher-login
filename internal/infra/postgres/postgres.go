package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/fx"
)

type PostgresParams struct {
	fx.In
	Env    *env.Environment
	Logger *slog.Logger
}

func New(params PostgresParams) (*bun.DB, error) {
	logger := params.Logger.With(
		slog.String("location", "postgres"),
		slog.String("layer", "infra"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := params.Env.GOPHER_POSTGRES_URL

	if dsn == "" {
		logger.ErrorContext(ctx, string(ports.CodeDatabaseConnFailed), "error", "DSN is empty")
		return nil, errors.New("database DSN is required")
	}

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	db := sql.OpenDB(connector)

	bunDB := bun.NewDB(db, pgdialect.New())
	bunDB.SetConnMaxIdleTime(5 * time.Minute)
	bunDB.SetMaxOpenConns(25)
	bunDB.SetMaxIdleConns(25)

	if err := bunDB.PingContext(ctx); err != nil {
		logger.ErrorContext(ctx, string(ports.CodeDatabasePingFailed), "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, string(ports.CodeDatabaseConnSuccess))
	return bunDB, nil
}
