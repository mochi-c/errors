package errors

type emptyInfo struct{}

func (emptyInfo) WhenError(cause error) string {
	if cause != nil {
		return cause.Error()
	} else {
		return ""
	}
}

func Wrap(err error) error {
	if err != nil {
		return withErrorInfo(err, emptyInfo{})
	} else {
		return nil
	}
}
