package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func New(env *env.Environment, logger *slog.Logger) (*bun.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := env.GOPHER_POSTGRES_URL
	fmt.Printf("dsn: %v\n", dsn)

	logbullHost := os.Getenv("GOPHER_LOGBULL_HOST")
	fmt.Printf("logbullHost: %v\n", logbullHost)
	if dsn == "" {
		logger.ErrorContext(ctx, ports.MsgDatabaseConnFailed, "error", "DSN is empty")
		return nil, errors.New("database DSN is required")
	}

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	db := sql.OpenDB(connector)

	bunDB := bun.NewDB(db, pgdialect.New())
	bunDB.SetConnMaxIdleTime(5 * time.Minute)
	bunDB.SetMaxOpenConns(25)
	bunDB.SetMaxIdleConns(25)

	if err := bunDB.PingContext(ctx); err != nil {
		logger.ErrorContext(ctx, ports.MsgDatabasePingFailed, "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, ports.MsgDatabaseConnSuccess)
	return bunDB, nil
}
