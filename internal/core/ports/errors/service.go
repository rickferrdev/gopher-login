package ports_errors

const (
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	CodeUserAlreadyExists  Code = "USER_ALREADY_EXISTS"
)

const (
	MessageInvalidCredentials Message = "invalid credentials"
	MessageUserAlreadyExists  Message = "user already exists"
)

func NewInvalidCredentials(trace error) error {
	return pick(CodeInvalidCredentials, trace)
}

func NewUserAlreadyExists(trace error) error {
	return pick(CodeUserAlreadyExists, trace)
}
