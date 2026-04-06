package ports

type Message string

var (
	// Resource Related
	MessageNotFound      Message = "resource not found"
	MessageAlreadyExists Message = "resource already exists"
	MessageInvalidID     Message = "provided identifier is invalid"
	MessageDataConflict  Message = "data conflict occurred"

	// Authentication & Security
	MessageInvalidCredentials Message = "invalid credentials"
	MessageUnauthorized       Message = "access denied"
	MessageInvalidToken       Message = "invalid or malformed token"
	MessageTokenExpired       Message = "session has expired"
	MessageSecurityError      Message = "security processing failed"

	// Request Validation
	MessageBadRequest       Message = "invalid request parameters"
	MessageInvalidInput     Message = "provided input is incorrect"
	MessageValidationFailed Message = "validation constraints not met"

	// System & Server
	MessageInternalError      Message = "an unexpected error occurred"
	MessageServiceUnavailable Message = "service is temporarily unavailable"
	MessageTimeout            Message = "the operation timed out"

	// Infrastructure & Persistence
	MessageStorageError    Message = "failed to process storage operation"
	MessageConnectionError Message = "could not establish connection to the service"
	MessageOperationFailed Message = "the requested operation could not be completed"
)

type Code string

const (
	CodeRequestBindingFailed Code = "request:binding_failed"
	CodeRequestInvalidID     Code = "request:invalid_id_format"

	CodeAuthInvalidCredentials Code = "auth:invalid_credentials"
	CodeAuthUnauthorized       Code = "auth:unauthorized"
	CodeAuthJwtInvalidFormat   Code = "auth:jwt_invalid_format"
	CodeAuthJwtVerifyFailed    Code = "auth:jwt_verification_failed"
	CodeAuthTokenValidated     Code = "auth:token_validated"
	CodeAuthLoginSuccess       Code = "auth:login_success"
	CodeAuthHashFailed         Code = "auth:password_hash_failed"
	CodeAuthTokenGenFailed     Code = "auth:token_generation_failed"

	CodeUserRegistered    Code = "user:registration_success"
	CodeUserAlreadyExists Code = "user:already_exists"
	CodeUserNotFound      Code = "user:not_found"
	CodeUserConflict      Code = "user:conflict"
	CodeUserFetchSuccess  Code = "user:fetch_success"

	CodeSystemInternalError Code = "system:internal_server_error"
	CodeSystemServiceFailed Code = "system:layer_service_failed"
	CodeSystemBadRequest    Code = "system:bad_request"

	CodeDatabaseCreateFailed Code = "db:create_failed"
	CodeDatabaseFetchFailed  Code = "db:fetch_failed"
	CodeDatabaseSchemaFailed Code = "db:schema_init_failed"
	CodeDatabasePingFailed   Code = "db:ping_failed"
	CodeDatabaseConnFailed   Code = "db:connection_failed"
	CodeDatabaseConnSuccess  Code = "db:connection_success"

	CodeServerShutdown Code = "server:shutting_down"
	CodeServerFailed   Code = "server:startup_failed"
	CodeServerStart    Code = "server:start_attempt"
)

type GopherError struct {
	Code    Code
	Message Message
	Err     error
	Status  int
}

func (e *GopherError) Error() string {
	if e.Err != nil {
		return string(e.Code) + ": " + e.Err.Error()
	}

	return string(e.Code) + ": " + string(e.Message)
}

func (e *GopherError) Unwrap() error {
	return e.Err
}

func (e *GopherError) Is(target error) bool {
	t, ok := target.(*GopherError)
	if !ok {
		return false
	}

	return e.Code == t.Code
}

func NewError(code Code, message Message, status int, err error) *GopherError {
	return &GopherError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}
