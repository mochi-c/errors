package errors

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(unwraper)
		if !ok {
			break
		}
		err = cause.Unwrap()
		if err == nil {
			return cause.(error)
		}
	}
	return err
}

func GetErrorInfo[T any](err error) (res T, ok bool) {
	for err != nil {
		if ins, hit := err.(Fundamental[T]); hit {
			return ins.GetErrorInfo(), true
		}
		cause, unwrap := err.(unwraper)
		if !unwrap {
			break
		}
		err = cause.Unwrap()
	}
	return res, ok
}

// GetAllErrorInfo latest with base at the first index
func GetAllErrorInfo[T any](err error) []T {
	res := make([]T, 0)
	for err != nil {
		if ins, hit := err.(Fundamental[T]); hit {
			res = append(res, ins.GetErrorInfo())
		}
		cause, unwrap := err.(unwraper)
		if !unwrap {
			break
		}
		err = cause.Unwrap()
	}
	return res
}

func GetOriginalErrorInfo[T any](err error) (res T, ok bool) {
	allInfo := GetAllErrorInfo[T](err)
	if len(allInfo) > 0 {
		return allInfo[len(allInfo)-1], true
	} else {
		return res, false
	}
}

// GetStack find the deepest Stack
func GetStack(err error) (Stack, bool) {
	for err != nil {
		if ins, hit := err.(HasStack); hit {
			return ins.GetStack(), true
		}
		cause, unwrap := err.(unwraper)
		if !unwrap {
			break
		}
		err = cause.Unwrap()
	}
	return nil, false
}

func GetStackCause(err error) (frame Frame, ok bool) {
	s, ok := GetStack(err)
	if ok {
		tracks := s.StackTrace()
		if len(tracks) >= 1 {
			return tracks[0], true
		}
	}
	return frame, false
}
