package apperror

import (
	"errors"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrEmailNotVerified  = errors.New("email not verified")
	ErrLinkExpired       = errors.New("link expired")
)
