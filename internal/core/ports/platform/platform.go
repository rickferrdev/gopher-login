//go:generate mockgen -source platform.go -destination ../../../tests/mocks/ports/platform.go -package ports_mocks
package ports_platform

import "time"

type UserClaims struct {
	UserID      string
	ValidatedAt time.Time
}

type Tokenizer interface {
	GenerateUserToken(id string) (string, error)
	ValidateUserToken(token string) (*UserClaims, error)
}

type Hasher interface {
	Generate(pass []byte) ([]byte, error)
	Validate(hash, pass []byte) error
}

type StructValidator interface {
	Validate(out any) error
}
