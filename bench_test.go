package errors

import (
	goError "errors"
	"fmt"
	"testing"
)

func BenchmarkGoErrorNew(b *testing.B) {
	goError.New("error")
}

func BenchmarkNew(b *testing.B) {
	New("error")
}

func BenchmarkErrorf(b *testing.B) {
	err := goError.New("error")
	b.ResetTimer()
	fmt.Errorf("wrap it: %w", err)
}

func BenchmarkWrap(b *testing.B) {
	err := goError.New("error")
	b.ResetTimer()
	Wrap(err)
}

func BenchmarkWithMessageAddStack(b *testing.B) {
	err := goError.New("error")
	b.ResetTimer()
	WithMessage(err, "message")
}

func BenchmarkWithMessage(b *testing.B) {
	err := New("error")
	b.ResetTimer()
	WithMessage(err, "message")
}

func BenchmarkMultiWithMessage(b *testing.B) {
	err := goError.New("error")
	err = WithMessage(err, "message")
	err = WithMessage(err, "message")
	err = WithMessage(err, "message")
	err = WithMessage(err, "message")
	b.ResetTimer()
	WithMessage(err, "message")
}

func BenchmarkMultiWrap(b *testing.B) {
	err := goError.New("error")
	err = Wrap(err)
	b.ResetTimer()
	Wrap(err)
}
