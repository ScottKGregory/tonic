package errors

import "fmt"

type GenericErr struct {
	Err error
}

func NewGenericError(err ...error) *GenericErr {
	if len(err) > 0 {
		return &GenericErr{Err: err[0]}
	}
	return &GenericErr{}
}

func (e *GenericErr) Unwrap() error { return e.Err }

func (e *GenericErr) Error() string {
	return fmt.Errorf("generic error: %w", e.Err).Error()
}

func (e *GenericErr) Is(err error) bool {
	_, ok := err.(*GenericErr)
	return ok
}

func (e *GenericErr) External() string {
	return "unexpected server error"
}
