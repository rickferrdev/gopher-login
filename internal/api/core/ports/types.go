package ports

import "time"

type ConsumerPayload struct {
	Nickname string
	Username string
}

type LoginInput struct {
	Email    string
	Password string
}

type RegisterInput struct {
	Nickname string
	Username string

	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type RegisterOutput struct {
	ID string
}

type TotokenClaims struct {
	ID string
}

type Totoken interface {
	GenerateToken(string) (string, error)
	VerifyToken(string) (*TotokenClaims, error)
}

type Validator interface {
	Validate(out any) error
}

const Timeout = 5 * time.Second
