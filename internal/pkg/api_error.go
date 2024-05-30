package pkg

import (
	"errors"
)

type APIError struct {
	APICode     `json:"api_code"`
	ErrMessages []string `json:"err_messages"`
	Err         error    `json:"-"` // do not export to client
}

func NewAPIError(code APICode, err error, msgs ...string) *APIError {

	if code == "" {
		code = Unknown
	}

	if err == nil {
		err = errors.New("no underlying error")
	}

	if len(msgs) == 0 {
		msgs = []string{}
	}

	return &APIError{
		APICode:     code,
		ErrMessages: msgs,
		Err:         err,
	}
}

var _ error = (*APIError)(nil)

func (e *APIError) Error() string {
	return e.Err.Error()
}

func (e *APIError) Unwrap() error {
	return e.Err
}
