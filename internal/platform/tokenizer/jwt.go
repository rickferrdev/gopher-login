package tokenizer

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rickferrdev/gopher-login/internal/config/env"
	ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"
	platform "github.com/rickferrdev/gopher-login/internal/core/ports/platform"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Platform struct {
	secret []byte
	env    *env.Env
}

type Params struct {
	fx.In

	Env *env.Env
}

func New(params Params) (platform.Tokenizer, error) {
	return &Platform{secret: []byte(params.Env.JwtSecret), env: params.Env}, nil
}

func (platform *Platform) GenerateUserToken(id string) (string, error) {
	duration, err := time.ParseDuration(platform.env.JwtExpiresIn)
	if err != nil {
		return "", ports_errors.NewJwtGenerateFailed(err)
	}

	claims := struct {
		UserID string
		jwt.RegisteredClaims
	}{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := newToken.SignedString(platform.secret)
	if err != nil {
		return "", ports_errors.NewJwtGenerateFailed(err)
	}

	return token, nil
}

func (p *Platform) ValidateUserToken(tokenString string) (*platform.UserClaims, error) {
	claims := struct {
		UserID string
		jwt.RegisteredClaims
	}{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return p.secret, nil
	})

	if err != nil {
		return nil, ports_errors.NewJwtTokenInvalid(err)
	}

	if !token.Valid {
		return nil, ports_errors.NewJwtTokenInvalid(err)
	}

	if claims.UserID == "" {
		return nil, ports_errors.NewJwtClaimsInvalid(nil)
	}

	return &platform.UserClaims{
		UserID: claims.UserID,
	}, nil
}
