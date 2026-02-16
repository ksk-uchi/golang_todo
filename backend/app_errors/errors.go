package app_errors

import "errors"

var (
	ErrTodoAlreadyDone = errors.New("cannot update a completed todo")
)
