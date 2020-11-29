package eth2api

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
