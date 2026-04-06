package consumer

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rickferrdev/gopher-login/internal/api/core/domain"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/fx"
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
			return nil, ports.NewError(ports.CodeRequestInvalidID, ports.MessageInvalidID, 400, err)
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

type RepositoryParams struct {
	fx.In
	BunDB *bun.DB
}

func New(params RepositoryParams) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := params.BunDB.NewCreateTable().
		Model((*Schema)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		return nil, ports.NewError(ports.CodeDatabaseSchemaFailed, ports.MessageStorageError, 500, err)
	}

	return &Repository{
		database: params.BunDB,
	}, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("email = ?", email).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, 404, err)
		}
		return nil, ports.NewError(ports.CodeDatabaseFetchFailed, ports.MessageInternalError, 500, err)
	}
	return schema.ToDomain(), nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, 404, err)
		}
		return nil, ports.NewError(ports.CodeDatabaseFetchFailed, ports.MessageInternalError, 500, err)
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
			return "", ports.NewError(ports.CodeUserAlreadyExists, ports.MessageAlreadyExists, 409, err)
		}
		return "", ports.NewError(ports.CodeDatabaseCreateFailed, ports.MessageStorageError, 500, err)
	}

	return id, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*domain.Consumer, error) {
	var schema Schema
	if err := r.database.NewSelect().Model(&schema).Where("username = ?", username).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.NewError(ports.CodeUserNotFound, ports.MessageNotFound, 404, err)
		}
		return nil, ports.NewError(ports.CodeDatabaseFetchFailed, ports.MessageInternalError, 500, err)
	}
	return schema.ToDomain(), nil
}
