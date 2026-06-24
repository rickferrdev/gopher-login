package hasher

import (
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	ports_platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var Provide = fx.Provide(New)

type Platform struct{}

func New() (ports_platform.Hasher, error) {
	return &Platform{}, nil
}

func (platform *Platform) Generate(pass []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, ports_errors.NewHasherGenerateFailed(err)
	}

	return hash, nil
}

func (platform *Platform) Validate(hash, pass []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, pass); err != nil {
		return ports_errors.NewHasherValidateFailed(err)
	}

	return nil
}
