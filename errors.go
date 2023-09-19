package gma

import "errors"

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidPath      = errors.New("invalid path")
)
