package httperrors

type HTTPStatusError interface {
	error
	StatusCode() int
}

type HTTPRequestError struct {
	status int
	err    error
}

func (e *HTTPRequestError) Error() string {
	return e.err.Error()
}

func (e *HTTPRequestError) Unwrap() error {
	return e.err
}

func (e *HTTPRequestError) StatusCode() int {
	return e.status
}

func NewHTTPRequestError(code int, err error) error {
	return &HTTPRequestError{
		status: code,
		err:    err,
	}
}
