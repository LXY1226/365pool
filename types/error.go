package types

type MyError struct {
	ErrMsg string
}

func (e *MyError) Error() string {
	return e.ErrMsg
}
