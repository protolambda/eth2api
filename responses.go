package eth2api

type Headers map[string]string

type PreparedResponse interface {
	//
	Code() uint
	// body to encode as response, may be nil
	Body() interface{}
	// headers to put into the response, may be nil
	Headers() Headers
}

func RespondBadInput(err error) PreparedResponse {
	// TODO
	return nil
}

func RespondNotFound(msg string) PreparedResponse {
	// TODO
	return nil
}

func RespondAccepted(msg string) PreparedResponse {
	return nil
}

func RespondOK(data interface{}) PreparedResponse {
	// TODO
	return nil
}

func RespondOKMsg(msg string) PreparedResponse {
	// TODO
	return nil
}

func RespondInternalError(err error) PreparedResponse {
	// TODO
	return nil
}

func RespondSyncing(msg string) PreparedResponse {
	// TODO
	return nil
}

type ApiError interface {
	error
	Code() uint
}

// ErrorResponse is a base struct with common error contents for responses.
type ErrorResponse struct {
	// The response solely consists of a ErrorMessage
	ErrorMessage
}

// Invalid request syntax.
type InvalidRequest struct {
	ErrorResponse
}

// Beacon node internal error.
type InternalError struct {
	ErrorResponse
}

// Beacon node is currently syncing, try again later.
type CurrentlySyncing struct {
	ErrorResponse
}
