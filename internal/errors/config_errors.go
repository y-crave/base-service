package errors

import "errors"

var (
	ErrUnmarshal     = errors.New("failed to unmarshal configuration")
	ErrValidRequired = errors.New("failed to validate configuration")
)
