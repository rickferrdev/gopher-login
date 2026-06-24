package ports_errors

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	Code    string
	Message string
)

type ErrorDefinition struct {
	Message Message
	Status  int
}

type Error struct {
	Code    Code    `json:"code"`
	Trace   error   `json:"-"`
	Status  int     `json:"status"`
	Message Message `json:"message"`
}

const (
	CodeStartupFailed Code = "STARTUP_FAILED"
)

const (
	MessageStartupFailed Message = "startup failed"
)

func NewStartupFailed(trace error) error {
	return pick(CodeStartupFailed, trace)
}

var definitions = map[Code]ErrorDefinition{
	CodeHttpBadRequest: {
		Message: MessageHttpBadRequest,
		Status:  http.StatusBadRequest,
	},
	CodeHttpInternalServer: {
		Message: MessageHttpInternalServer,
		Status:  http.StatusInternalServerError,
	},
	CodeHttpUnauthorized: {
		Message: MessageHttpUnauthorized,
		Status:  http.StatusUnauthorized,
	},

	CodeInvalidCredentials: {
		Message: MessageInvalidCredentials,
		Status:  http.StatusUnauthorized,
	},
	CodeUserAlreadyExists: {
		Message: MessageUserAlreadyExists,
		Status:  http.StatusConflict,
	},

	CodeDatabaseNotFound: {
		Message: MessageDatabaseNotFound,
		Status:  http.StatusNotFound,
	},
	CodeDatabaseDuplicateKey: {
		Message: MessageDatabaseDuplicateKey,
		Status:  http.StatusConflict,
	},
	CodeDatabaseFailedConnect: {
		Message: MessageDatabaseFailedConnect,
		Status:  http.StatusInternalServerError,
	},
	CodeDatabaseInvalidID: {
		Message: MessageDatabaseInvalidID,
		Status:  http.StatusBadRequest,
	},
	CodeDatabaseInternal: {
		Message: MessageDatabaseInternal,
		Status:  http.StatusInternalServerError,
	},
	CodeDatabaseSchemaFailed: {
		Message: MessageDatabaseSchemaFailed,
		Status:  http.StatusInternalServerError,
	},

	CodeAMQPFailedConnect: {
		Message: MessageAMQPFailedConnect,
		Status:  http.StatusInternalServerError,
	},
	CodeAMQPFailedPublish: {
		Message: MessageAMQPFailedPublish,
		Status:  http.StatusInternalServerError,
	},
	CodeAMQPPayloadInvalid: {
		Message: MessageAMQPPayloadInvalid,
		Status:  http.StatusBadRequest,
	},
	CodeAMQPInternal: {
		Message: MessageAMQPInternal,
		Status:  http.StatusInternalServerError,
	},

	CodeHasherGenerateFailed: {
		Message: MessageHasherGenerateFailed,
		Status:  http.StatusInternalServerError,
	},
	CodeHasherValidateFailed: {
		Message: MessageHasherValidateFailed,
		Status:  http.StatusInternalServerError,
	},

	CodeMailerSendFailed: {
		Message: MessageMailerSendFailed,
		Status:  http.StatusInternalServerError,
	},

	CodeJwtGenerateFailed: {
		Message: MessageJwtGenerateFailed,
		Status:  http.StatusInternalServerError,
	},
	CodeJwtTokenInvalid: {
		Message: MessageJwtTokenInvalid,
		Status:  http.StatusUnauthorized,
	},
	CodeJwtClaimsInvalid: {
		Message: MessageJwtClaimsInvalid,
		Status:  http.StatusUnauthorized,
	},
	CodeJwtTokenExpired: {
		Message: MessageJwtTokenExpired,
		Status:  http.StatusUnauthorized,
	},

	CodeStartupFailed: {
		Message: MessageStartupFailed,
		Status:  http.StatusInternalServerError,
	},
}

func pick(code Code, trace error) error {
	value, ok := definitions[code]
	if !ok {
		return &Error{
			Code:    CodeHttpInternalServer,
			Message: MessageHttpInternalServer,
			Status:  http.StatusInternalServerError,
			Trace:   trace,
		}
	}

	return &Error{
		Code:    code,
		Message: value.Message,
		Status:  value.Status,
		Trace:   trace,
	}
}

func (e *Error) Error() string {
	if e.Trace != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.Code, e.Status, e.Message, e.Trace)
	}

	return fmt.Sprintf("[%s:%d] %s", e.Code, e.Status, e.Message)
}

func (e *Error) Unwrap() error {
	if e == nil || e.Trace == nil {
		return nil
	}

	return e.Trace
}

func IsCode(err error, code Code) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

func NewError(code Code, message Message, status int, trace error) error {
	return &Error{
		Code:    code,
		Trace:   trace,
		Status:  status,
		Message: message,
	}
}
