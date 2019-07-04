package types

type MyError struct {
	errMsg string
}

func (e *MyError) Error() string {
	return e.errMsg
}
