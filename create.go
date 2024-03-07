package errors

import (
	"errors"
	"fmt"
)

func New(msg string) error {
	cause := errors.New(msg)
	return withErrorInfo(cause, emptyInfo{})
}

func Errorf(format string, args ...interface{}) error {
	cause := fmt.Errorf(format, args...)
	return withErrorInfo(cause, emptyInfo{})
}
