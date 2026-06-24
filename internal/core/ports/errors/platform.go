package ports_errors

const (
	CodeAMQPFailedConnect  Code = "AMQP_FAILED_CONNECT"
	CodeAMQPFailedPublish  Code = "AMQP_FAILED_PUBLISH"
	CodeAMQPPayloadInvalid Code = "AMQP_PAYLOAD_INVALID"
	CodeAMQPInternal       Code = "AMQP_INTERNAL"

	CodeHasherGenerateFailed Code = "HASHER_GENERATE_FAILED"
	CodeHasherValidateFailed Code = "HASHER_VALIDATE_FAILED"

	CodeMailerSendFailed Code = "MAILER_SEND_FAILED"

	CodeJwtGenerateFailed Code = "JWT_GENERATE_FAILED"
	CodeJwtTokenInvalid   Code = "JWT_TOKEN_INVALID"
	CodeJwtClaimsInvalid  Code = "JWT_CLAIMS_INVALID"
	CodeJwtTokenExpired   Code = "JWT_CLAIMS_EXPIRED"
)

const (
	MessageAMQPFailedConnect  Message = "amqp failed to connect"
	MessageAMQPFailedPublish  Message = "amqp failed to publish"
	MessageAMQPPayloadInvalid Message = "amqp payload invalid"
	MessageAMQPInternal       Message = "amqp internal error"

	MessageHasherGenerateFailed Message = "hasher generate failed"
	MessageHasherValidateFailed Message = "hasher validate failed"

	MessageMailerSendFailed Message = "mailer send failed"

	MessageJwtGenerateFailed Message = "jwt generate failed"
	MessageJwtTokenInvalid   Message = "jwt token invalid"
	MessageJwtClaimsInvalid  Message = "jwt claims invalid"
	MessageJwtTokenExpired   Message = "jwt token expired"
)

func NewAMQPFailedConnect(trace error) error {
	return pick(CodeAMQPFailedConnect, trace)
}

func NewAMQPFailedPublish(trace error) error {
	return pick(CodeAMQPFailedPublish, trace)
}

func NewAMQPPayloadInvalid(trace error) error {
	return pick(CodeAMQPPayloadInvalid, trace)
}

func NewAMQPInternal(trace error) error {
	return pick(CodeAMQPInternal, trace)
}

func NewMailerSendFailed(trace error) error {
	return pick(CodeMailerSendFailed, trace)
}

func NewHasherGenerateFailed(trace error) error {
	return pick(CodeHasherGenerateFailed, trace)
}

func NewHasherValidateFailed(trace error) error {
	return pick(CodeHasherValidateFailed, trace)
}

func NewJwtGenerateFailed(trace error) error {
	return pick(CodeJwtGenerateFailed, trace)
}

func NewJwtTokenInvalid(trace error) error {
	return pick(CodeJwtTokenInvalid, trace)
}

func NewJwtClaimsInvalid(trace error) error {
	return pick(CodeJwtClaimsInvalid, trace)
}

func NewJwtTokenExpired(trace error) error {
	return pick(CodeJwtTokenExpired, trace)
}
