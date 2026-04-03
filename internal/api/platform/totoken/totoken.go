package totoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rickferrdev/gopher-login/internal/api/core/ports"
	"github.com/rickferrdev/gopher-login/internal/config/env"
)

var (
	ErrInvalidToken            = errors.New("invalid or expired token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrMissingSubject          = errors.New("token subject is missing")
)

type TotokenAdapter struct {
	secret []byte
}

func New(env *env.Environment) *TotokenAdapter {
	return &TotokenAdapter{
		secret: []byte(env.GOPHER_SERVER_JWT_SECRET),
	}
}

func (a *TotokenAdapter) GenerateToken(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

func (a *TotokenAdapter) VerifyToken(tokenString string) (*ports.TotokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return a.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject == "" {
			return nil, ErrMissingSubject
		}

		return &ports.TotokenClaims{
			ID: claims.Subject,
		}, nil
	}

	return nil, ErrInvalidToken
}
