package ports_errors

const (
	CodeHttpBadRequest     Code = "HTTP_BAD_REQUEST"
	CodeHttpInternalServer Code = "HTTP_INTERNAL_SERVER"
	CodeHttpUnauthorized   Code = "HTTP_UNAUTHORIZED"
)

const (
	MessageHttpBadRequest     Message = "bad request"
	MessageHttpInternalServer Message = "internal server error"
	MessageHttpUnauthorized   Message = "unauthorized"
)

func NewHttpBadRequest(trace error) error {
	return pick(CodeHttpBadRequest, trace)
}

func NewHttpInternalServer(trace error) error {
	return pick(CodeHttpInternalServer, trace)
}

func NewHttpUnauthorized(trace error) error {
	return pick(CodeHttpUnauthorized, trace)
}
