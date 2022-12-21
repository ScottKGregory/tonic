package errors

import "errors"

type TonicError interface {
	External() string
	Error() string
	Is(err error) bool
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
