package user

import (
	"time"

	"github.com/rickferrdev/gopher-login/internal/core/domain"
	"github.com/rickferrdev/gopher-login/internal/outbound/databases/mongo/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID       bson.ObjectID `bson:"_id"`
	Email    string        `bson:"email"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (user *User) ToDomain() (*domain.User, error) {
	toDomain := domain.User{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &toDomain, nil
}

func FromUserDomain(user domain.User) (*User, error) {
	var id bson.ObjectID
	if user.ID != "" {
		newFromID, err := bson.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, schema.NewSchemaInvalidObjectIDFromHex(err)
		}

		id = newFromID
	} else {
		return nil, schema.NewMissingInvalidObjectID(nil)
	}

	toSchema := User{
		ID:        id,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &toSchema, nil
}
