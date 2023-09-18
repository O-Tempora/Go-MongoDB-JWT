package resperr

import "errors"

var (
	ErrInvalidGUID        = errors.New("GUID input does not match any existing users")
	ErrInvalidToken       = errors.New("Failed to validate Access and Refresh token pair")
	ErrInvalidRequestBody = errors.New("Invalid request body")
)
