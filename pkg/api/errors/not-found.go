package errors

import (
	"fmt"

	"github.com/scottkgregory/tonic/pkg/helpers"
)

type NotFoundErr struct {
	ID string
}

func NewNotFoundError(id string) *NotFoundErr {
	return &NotFoundErr{ID: id}
}

func (e *NotFoundErr) Error() string {
	if helpers.IsEmptyOrWhitespace(e.ID) {
		return "not found"
	}

	return fmt.Sprintf("not found: %s", e.ID)
}

func (e *NotFoundErr) Is(err error) bool {
	_, ok := err.(*NotFoundErr)
	return ok
}
