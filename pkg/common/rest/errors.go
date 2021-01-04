package rest

import (
	"errors"
	"fmt"
)

type HandlerError struct {
	Code    int
	err     error
	message string
}

func (err *HandlerError) Error() string {
	if err.message == "" {
		return err.err.Error()
	}
	return fmt.Sprintf("%s: %s", err.message, err.err.Error())
}

func (err *HandlerError) Cause() error {
	return err.err
}

func WarpError(err error, code int, message string) error {
	return &HandlerError{
		Code:    code,
		err:     err,
		message: message,
	}
}

func NewError(code int, message string) error {
	return &HandlerError{
		Code:    code,
		err:     errors.New(message),
		message: "",
	}
}
