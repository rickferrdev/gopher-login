package ports_errors

const (
	CodeDatabaseNotFound      Code = "DATABASE_NOT_FOUND"
	CodeDatabaseDuplicateKey  Code = "DATABASE_DUPLICATE_KEY"
	CodeDatabaseFailedConnect Code = "DATABASE_FAILED_CONNECT"
	CodeDatabaseInvalidID     Code = "DATABASE_INVALID_ID"
	CodeDatabaseInternal      Code = "DATABASE_INTERNAL"
	CodeDatabaseSchemaFailed  Code = "DATABASE_SCHEMA_FAILED"
)

const (
	MessageDatabaseNotFound      Message = "database record not found"
	MessageDatabaseDuplicateKey  Message = "database duplicate key"
	MessageDatabaseFailedConnect Message = "database failed to connect"
	MessageDatabaseInvalidID     Message = "database invalid id"
	MessageDatabaseInternal      Message = "database internal error"
	MessageDatabaseSchemaFailed  Message = "database schema failed"
)

func NewDatabaseFailedConnect(trace error) error {
	return pick(CodeDatabaseFailedConnect, trace)
}

func NewDatabaseNotFound(trace error) error {
	return pick(CodeDatabaseNotFound, trace)
}

func NewDatabaseInvalidID(trace error) error {
	return pick(CodeDatabaseInvalidID, trace)
}

func NewDatabaseInternal(trace error) error {
	return pick(CodeDatabaseInternal, trace)
}

func NewDatabaseDuplicateKey(trace error) error {
	return pick(CodeDatabaseDuplicateKey, trace)
}

func NewDatabaseSchemaFailed(trace error) error {
	return pick(CodeDatabaseSchemaFailed, trace)
}
