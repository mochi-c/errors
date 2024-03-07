package errors

import "fmt"

type message string

func (m message) WhenError(cause error) string {
	if cause == nil {
		return string(m)
	} else {
		return string(m) + ": " + cause.Error()
	}
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return withErrorInfo(err, message(msg))
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(args) > 0 {
		return withErrorInfo(err, message(fmt.Sprintf(format, args...)))
	} else {
		return withErrorInfo(err, message(format))
	}
}
