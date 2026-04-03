package ports

import "errors"

var (
	ErrConsumerNotFound      = errors.New("consumer not found")
	ErrConsumerAlreadyExists = errors.New("consumer with given email or username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrConsumerInvalidID     = errors.New("invalid consumer ID format")
	ErrContextTimeout        = errors.New("context timeout")

	ErrInvalidToken   = errors.New("invalid token")
	ErrInternalServer = errors.New("internal server error")
)

const (
	// Contexto: Request & Binding
	MsgRequestBindingFailed = "request:binding_failed"
	MsgRequestInvalidID     = "request:invalid_id_format"

	// Contexto: Auth & JWT
	MsgAuthInvalidCredentials = "auth:invalid_credentials"
	MsgAuthUnauthorized       = "auth:unauthorized"
	MsgAuthJwtInvalidFormat   = "auth:jwt_invalid_format"
	MsgAuthJwtVerifyFailed    = "auth:jwt_verification_failed"
	MsgAuthTokenValidated     = "auth:token_validated"
	MsgAuthLoginSuccess       = "auth:login_success"
	MsgAuthHashFailed         = "auth:password_hash_failed"
	MsgAuthTokenGenFailed     = "auth:token_generation_failed"

	// Contexto: User Domain
	MsgUserRegistered    = "user:registration_success"
	MsgUserAlreadyExists = "user:already_exists"
	MsgUserNotFound      = "user:not_found"
	MsgUserConflict      = "user:conflict"
	MsgUserFetchSuccess  = "user:fetch_success"

	// Contexto: System & Infra
	MsgSystemInternalError = "system:internal_server_error"
	MsgSystemServiceFailed = "system:layer_service_failed"
	MsgSystemBadRequest    = "system:bad_request"

	MsgDatabaseCreateFailed = "db:create_failed"
	MsgDatabaseFetchFailed  = "db:fetch_failed"
	MsgDatabaseSchemaFailed = "db:schema_init_failed"
	MsgDatabasePingFailed   = "db:ping_failed"
	MsgDatabaseConnFailed   = "db:connection_failed"
	MsgDatabaseConnSuccess  = "db:connection_success"

	MsgServerShutdown = "server:shutting_down"
	MsgServerFailed   = "server:startup_failed"
	MsgServerStart    = "server:start_attempt"
)
