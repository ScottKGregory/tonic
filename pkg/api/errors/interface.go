package errors

type TonicError interface {
	External() string
	Error() string
	Is(err error) bool
}
