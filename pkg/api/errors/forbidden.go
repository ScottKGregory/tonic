package errors

type ForbiddenErr struct {
	Required []string
}

func NewForbiddenError(required ...string) *ForbiddenErr {
	return &ForbiddenErr{required}
}

func (e *ForbiddenErr) Error() string {
	return "forbidden"
}

func (e *ForbiddenErr) Is(err error) bool {
	_, ok := err.(*ForbiddenErr)
	return ok
}

func (e *ForbiddenErr) External() string {
	return "forbidden"
}
