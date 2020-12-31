package errors

type HandlerError struct {
	Code    int
	Message string
}

func (err HandlerError) Error() string {
	return err.Message
}
