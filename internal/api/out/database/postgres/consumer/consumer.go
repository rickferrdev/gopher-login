package consumer

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Schema struct {
	bun.BaseModel `bun:"table:consumers,alias:c"`

	ID       uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Username string    `bun:"username,unique,notnull"`
	Nickname string    `bun:"nickname,notnull"`
	Email    string    `bun:"email,unique,notnull"`
	Password string    `bun:"password,notnull"`
}

func NewSchema(consumer domain.Consumer) (*Schema, error) {
	schema := &Schema{
		Username: consumer.Username,
		Nickname: consumer.Nickname,
		Email:    consumer.Email,
		Password: consumer.Password,
	}

	if consumer.ID != "" {
		parsedID, err := uuid.Parse(consumer.ID)
		if err != nil {
			slog.Error(ports.MsgRequestInvalidID, "id", consumer.ID, "error", err)
			return nil, ports.ErrConsumerInvalidID
		}
		schema.ID = parsedID
	} else {
		schema.ID = uuid.New()
	}

	return schema, nil
}

func (s *Schema) ToDomain() *domain.Consumer {
	return &domain.Consumer{
		ID:       s.ID.String(),
		Username: s.Username,
		Nickname: s.Nickname,
		Email:    s.Email,
		Password: s.Password,
	}
}

type Repository struct {
	database *bun.DB
}

func New(bunDB *bun.DB) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := bunDB.NewCreateTable().
		Model((*Schema)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		slog.Error(ports.MsgDatabaseSchemaFailed, "error", err)
		return nil, err
	}

	return &Repository{
		database: bunDB,
	}, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("email = ?", email).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrConsumerNotFound
		}
		slog.ErrorContext(ctx, ports.MsgDatabaseFetchFailed, "email", email, "error", err)
		return nil, ports.ErrInternalServer
	}
	return schema.ToDomain(), nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrConsumerNotFound
		}
		slog.ErrorContext(ctx, ports.MsgDatabaseFetchFailed, "id", id, "error", err)
		return nil, ports.ErrInternalServer
	}
	return schema.ToDomain(), nil
}

func (r *Repository) Create(ctx context.Context, consumer domain.Consumer) (string, error) {
	schema, err := NewSchema(consumer)
	if err != nil {
		return "", err
	}

	var id string
	if _, err := r.database.NewInsert().Model(schema).Returning("id").Exec(ctx, &id); err != nil {
		if pgErr, ok := err.(pgdriver.Error); ok && pgErr.Field('C') == "23505" {
			slog.WarnContext(ctx, ports.MsgUserAlreadyExists, "id", schema.ID)
			return "", ports.ErrConsumerAlreadyExists
		}

		slog.ErrorContext(ctx, ports.MsgDatabaseCreateFailed, "id", schema.ID, "error", err)
		return "", ports.ErrInternalServer
	}

	return id, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("username = ?", username).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrConsumerNotFound
		}
		slog.ErrorContext(ctx, ports.MsgDatabaseFetchFailed, "username", username, "error", err)
		return nil, ports.ErrInternalServer
	}

	return schema.ToDomain(), nil
}
