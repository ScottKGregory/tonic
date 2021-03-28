package errors

type UnauthorisedErr struct {
}

func NewUnauthorisedError() *UnauthorisedErr {
	return &UnauthorisedErr{}
}

func (e *UnauthorisedErr) Error() string {
	return "unauthorised"
}

func (e *UnauthorisedErr) Is(err error) bool {
	_, ok := err.(*UnauthorisedErr)
	return ok
}
