package web

// ErrorResponse how we respond to clients when something goes wrong.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to add web information to a request error.
type Error struct {
	Err  error
	Code int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewRequestError(err error, code int) error {
	return &Error{Err: err, Code: code}
}
