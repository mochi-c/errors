package errors

import (
	goError "errors"
	"fmt"
)

func ExampleWithMessage() {
	err := goError.New("whoops")
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithMessage(err, "msg3")
	fmt.Println(err)

	// Output: msg3: msg2: msg1: whoops
}

func ExampleCause() {
	err := goError.New("whoops")
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithMessage(err, "msg3")
	fmt.Println(Cause(err))

	// Output: whoops
}

func ExampleGetErrorInfo() {
	err := WithErrorInfo(New("New"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	baseInfo, ok := GetErrorInfo[CodeInfo](err)
	fmt.Println(baseInfo.Code, ok)

	// Output: 101 true
}

func ExampleGetOriginalErrorInfo() {
	err := WithErrorInfo(New("New"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	baseInfo, ok := GetOriginalErrorInfo[CodeInfo](err)
	fmt.Println(baseInfo.Code, ok)

	// Output: 100 true
}

func ExampleGetAllErrorInfo() {
	err := WithErrorInfo(New("msg"), CodeInfo{
		Code: 100,
	})
	err = WithMessage(err, "msg1")
	err = WithMessage(err, "msg2")
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	baseInfos := GetAllErrorInfo[CodeInfo](err)
	fmt.Println(baseInfos[0].Code)
	fmt.Println(baseInfos[1].Code)

	// Output:
	//101
	//100
}

func ExampleGetStack() {
	err := New("whoops")
	stack, _ := GetStack(err)
	fmt.Printf("%v", stack)

	/*
		like this

		code.byted.org/tns/algo_eng_errors.ExampleGetStack
			/Users/bytedance/workspace/tiktok/errors/example_test.go:80
		testing.runExample
			/Users/bytedance/sdk/go1.21.4/src/testing/run_example.go:63
		testing.runExamples
			/Users/bytedance/sdk/go1.21.4/src/testing/example.go:44
		testing.(*M).Run
			/Users/bytedance/sdk/go1.21.4/src/testing/testing.go:1927
		main.main
			_testmain.go:97
		runtime.main
			/Users/bytedance/sdk/go1.21.4/src/runtime/proc.go:267
		runtime.goexit
			/Users/bytedance/sdk/go1.21.4/src/runtime/asm_arm64.s:1197
	*/

}

func ExampleGetStackCause() {
	err := New("whoops")
	frame, ok := GetStackCause(err)
	fmt.Println(frame.FuncName(), frame.Line(), ok)

	// Output:	ExampleGetStackCause 106 true

	/*
		fmt.Println(frame.FuncName(), frame.File(), frame.Line(), ok)
		like this
		ExampleGetStackCause /Users/bytedance/workspace/tiktok/errors/example_test.go 106 true
	*/

}
