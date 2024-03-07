
<h1 >ErrorWithStack</h1>

Quickly construct an Error carrying additional context information and automatically pass or append stack information.

## Example

```golang
package main

func ExampleHowToUse() {

	var err error

	//Create Error
	err = New("whoops")

	//Warp some error have no stack
	var originalError error
	err = Wrap(originalError)

        //add some message and stack
	err = WithMessage(err, "msg1")
	
	//add some ErrorCode
	err = WithErrorInfo(err, CodeInfo{
		Code: 101,
	})

	//get errorcode
	codeInfo, ok := GetOriginalErrorInfo[CodeInfo](err)
	fmt.Println(codeInfo.Code, ok)

	//get stack
	stack, ok := GetStack(err)
	fmt.Printf("%v", stack)

	//get cause stack location
	cause, ok := GetStackCause(err)
	fmt.Printf("%+v", cause)

	/*
		see more example_test.go
	*/

}

```

## Suggestions for using Error

All errors should carry stack and cause information and allow for easy appending of any context data

In principle, at any time, either handle the error (usually by restoring the scene or switching to fallback logic and recording the context) and stop throwing the error, or append the current context info and continue throwing the error. A typical counterexample is throwing the error while printing the error log.

When throw the error, use **WithMessage/WithMessagef/WithErrorInfo/Wrap** to append context information, or just append the stack. When Cause is nil, those functions will directly return nil. For any cause without stack information, these functions will append the current stack information. If the error already has stack information, it will not be appended repeatedly. Therefore, they can be easily and universally used.

You can easily create a new Error with stack information using **New** or **Errorf**.