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

type BasicResponse struct {
	code    uint
	body    any
	headers Headers
}

func (b *BasicResponse) Code() uint {
	return b.code
}

func (b *BasicResponse) Body() interface{} {
	return b.body
}

func (b *BasicResponse) Headers() Headers {
	return b.headers
}

var _ PreparedResponse = (*BasicResponse)(nil)

func RespondBadInputs(msg string, failures []IndexedErrorMessageItem) PreparedResponse {
	return &BasicResponse{
		code: 400,
		body: &ErrorMessage{
			CodeValue: 400,
			Message:   msg,
			Failures:  failures,
		},
		headers: nil,
	}
}

func RespondBadInput(err error) PreparedResponse {
	return &BasicResponse{
		code: 400,
		body: ErrorMessage{CodeValue: 400, Message: err.Error()},
	}
}

func RespondNotFound(msg string) PreparedResponse {
	return &BasicResponse{
		code: 404,
		body: &ErrorMessage{
			CodeValue: 404,
			Message:   msg,
		},
	}
}

func RespondAccepted(err error) PreparedResponse {
	return &BasicResponse{
		code: 202,
		body: ErrorMessage{CodeValue: 202, Message: err.Error()},
	}
}

func RespondOK(body interface{}) PreparedResponse {
	return &BasicResponse{
		code: 200,
		body: body,
	}
}

func RespondOKMsg(msg string) PreparedResponse {
	return &BasicResponse{
		code: 200,
		body: &ErrorMessage{
			CodeValue: 200,
			Message:   msg,
		},
	}
}

func RespondInternalError(err error) PreparedResponse {
	return &BasicResponse{
		code: 500,
		body: &ErrorMessage{
			CodeValue: 500,
			Message:   err.Error(),
		},
	}
}

func RespondSyncing(msg string) PreparedResponse {
	return &BasicResponse{
		code: 503,
		body: &ErrorMessage{
			CodeValue: 503,
			Message:   msg,
		},
	}
}

type ApiError interface {
	error
	Code() uint
}
