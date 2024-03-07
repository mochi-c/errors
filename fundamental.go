package errors

type ErrorInfo interface {
	WhenError(cause error) string
}

type unwraper interface {
	Unwrap() error
}

type HasStack interface {
	GetStack() Stack
}

type Fundamental[T any] interface {
	HasStack
	GetErrorInfo() T
}

type fundamental[T ErrorInfo] struct {
	cause error
	info  T
	stack Stack
}

func (f *fundamental[T]) GetErrorInfo() T {
	return f.info
}

func (f *fundamental[T]) GetStack() Stack {
	return f.stack
}

func (f *fundamental[T]) Error() string {
	return f.info.WhenError(f.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (f *fundamental[T]) Unwrap() error {
	return f.cause
}

func WithErrorInfo[T ErrorInfo](err error, info T) error {
	return withErrorInfo(err, info)
}

func withErrorInfo[T ErrorInfo](err error, info T) error {
	if err == nil {
		return nil
	}
	if stack, ok := GetStack(err); ok {
		return &fundamental[T]{
			cause: err,
			info:  info,
			stack: stack,
		}
	} else {
		stack = callers(4)
		return &fundamental[T]{
			cause: err,
			info:  info,
			stack: stack,
		}
	}

}

func newFundamental[T ErrorInfo](info T) error {
	return &fundamental[T]{
		cause: nil,
		info:  info,
		stack: callers(4),
	}
}
