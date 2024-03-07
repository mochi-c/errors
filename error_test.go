package errors

import (
	goError "errors"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("%s %s %s", "foo", "foo", "foo"), "foo foo foo"},
		{New("foo"), "foo"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.want {
			t.Errorf("got: %q, want %q", tt.err.Error(), tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := goError.New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  WithMessage(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		WithMessage(nil, "whoops"),
		nil,
	}, {
		WithMessage(io.EOF, "whoops"),
		io.EOF,
	}, {
		Wrap(nil),
		nil,
	}, {
		Wrap(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil)
	if got != nil {
		t.Errorf("Wrapf(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{Wrap(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err).Error()
		if got != tt.want {
			t.Errorf("WithStack(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := WithMessage(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := WithMessagef(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		goError.New("EOF"),
		New("EOF"),
		WithMessage(io.EOF, "EOF"),
		WithMessagef(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		Wrap(io.EOF),
		Wrap(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

type CodeInfo struct {
	Code int
}

func (info CodeInfo) WhenError(cause error) string {
	if cause != nil {
		return cause.Error()
	} else {
		return fmt.Sprintf("%d", info.Code)
	}
}

func TestGetErrorInfo(t *testing.T) {
	err := WithErrorInfo(New("New"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	baseInfo, ok := GetErrorInfo[CodeInfo](err)
	if ok != true || baseInfo.Code != 101 {
		t.Fail()
	}
}

func TestExampleGetAllErrorInfo(t *testing.T) {
	err := WithErrorInfo(New("msg"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	baseInfos := GetAllErrorInfo[CodeInfo](err)
	if baseInfos[0].Code != 101 || baseInfos[1].Code != 100 {
		t.Fail()
	}
}

func TestExampleGetOriginalErrorInfo(t *testing.T) {
	err := WithErrorInfo(New("msg"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	info, ok := GetOriginalErrorInfo[CodeInfo](err)
	if ok != true || info.Code != 100 {
		t.Fail()
	}
}

func TestStack(t *testing.T) {
	goErr := goError.New("whoops")
	goErrWarp := Wrap(goErr)
	goErrWithMessage := WithMessage(goErr, "message")
	goErrWithMessagef := WithMessagef(goErr, "message %s", "message")
	goErrWithInfo := WithErrorInfo(goErr, CodeInfo{
		Code: 100,
	})

	newError := New("whoops")
	multiWarp := Wrap(newError)
	multiWarp = Wrap(newError)
	withMessage := WithMessage(newError, "message")

	tests := []struct {
		name     string
		err      error
		funcName string
		line     int
	}{
		{"goErrWarp", goErrWarp, "TestStack", 243},
		{"goErrWithMessage", goErrWithMessage, "TestStack", 244},
		{"goErrWithMessagef", goErrWithMessagef, "TestStack", 245},
		{"goErrWithInfo", goErrWithInfo, "TestStack", 246},
		{"err", newError, "TestStack", 250},
		{"multiWarp", multiWarp, "TestStack", 250},
		{"withMessage", withMessage, "TestStack", 250},
	}

	for _, tt := range tests {
		stack, ok := GetStack(tt.err)
		if !ok {
			t.Errorf("get stack fail %v", tt.name)
		}
		if stack.StackSource().FuncName() != tt.funcName || stack.StackSource().Line() != tt.line {
			t.Errorf("%v funcName need %s but is %s line need %d but is %d", tt.name, tt.funcName, stack.StackSource().FuncName(), tt.line, stack.StackSource().Line())
		}

		stackCause, ok := GetStackCause(tt.err)
		if !ok {
			t.Errorf("get stack cause fail %v", tt.name)
		}
		if stackCause.FuncName() != tt.funcName || stackCause.Line() != tt.line {
			t.Errorf("%v funcName need %s but is %s line need %d but is %d", tt.name, tt.funcName, stackCause.FuncName(), tt.line, stackCause.Line())
		}
	}
}

func TestStackFormat(t *testing.T) {
	newError := New("whoops")
	multiWarp := Wrap(newError)
	multiWarp = Wrap(newError)

	stack, _ := GetStack(multiWarp)

	fmt.Printf("%v", stack)

	stackCause, _ := GetStackCause(multiWarp)
	fmt.Printf("%+v\n", stackCause)

}
