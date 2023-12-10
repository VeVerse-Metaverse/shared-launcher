package errors

type WrappedError struct {
	Message string
	Err     error
}

func (e WrappedError) Error() string {
	return e.Message
}

func (e WrappedError) Unwrap() error {
	return e.Err
}

func (e WrappedError) Is(target error) bool {
	err, ok := target.(WrappedError)
	if ok {
		return err.Message == e.Message
	}
	return ok
}
