package errors

type ForbiddenErr struct {
}

func NewForbiddenError() *ForbiddenErr {
	return &ForbiddenErr{}
}

func (e *ForbiddenErr) Error() string {
	return "forbidden"
}

func (e *ForbiddenErr) Is(err error) bool {
	_, ok := err.(*ForbiddenErr)
	return ok
}
