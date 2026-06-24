package schema

import ports_errors "github.com/rickferrdev/gopher-login/internal/core/ports/errors"

const (
	CodeSchemaInvalidObjectIDFromHex ports_errors.Code = "DATABASE_INVALID_OBJECT_ID_FROM_HEX"
	CodeMissingValidObjectID         ports_errors.Code = "MISSING_INVALID_OBJECT_ID"
)

const (
	MessageMissingValidObjectID         ports_errors.Message = "ID is expected from this object"
	MessageSchemaInvalidObjectIDFromHex ports_errors.Message = "received ID returned an error when converted to the bson.ObjectID type"
)

func NewSchemaInvalidObjectIDFromHex(trace error) error {
	return ports_errors.NewError(
		CodeSchemaInvalidObjectIDFromHex,
		MessageSchemaInvalidObjectIDFromHex,
		500,
		trace,
	)
}

func NewMissingInvalidObjectID(trace error) error {
	return ports_errors.NewError(
		CodeMissingValidObjectID,
		MessageMissingValidObjectID,
		500,
		nil,
	)
}
