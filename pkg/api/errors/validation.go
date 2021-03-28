package errors

import (
	"errors"
	"fmt"

	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
)

var InvalidErr = errors.New("invalid")

type ValidationErr struct {
	Message    string
	Validation map[string]string
}

func NewValidationError(validation ...map[string]string) *ValidationErr {
	var val map[string]string
	if len(validation) > 0 {
		val = validation[0]
	}

	err := &ValidationErr{}

	if m, ok := val[constants.GlobalKey]; ok {
		err.Message = m
		delete(val, constants.GlobalKey)
	}

	err.Validation = val

	return err
}

func (e *ValidationErr) Error() string {
	if !helpers.IsEmptyOrWhitespace(e.Message) {
		return fmt.Sprintf("validation failed: %s", e.Message)
	}

	return "validation failed"
}

func (e *ValidationErr) Is(err error) bool {
	_, ok := err.(*ValidationErr)
	return ok
}
