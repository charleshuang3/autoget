package errors

type HTTPStatusError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHTTPStatusError(code int, message string) *HTTPStatusError {
	return &HTTPStatusError{
		Code:    code,
		Message: message,
	}
}

func (e *HTTPStatusError) Error() string {
	return e.Message
}
