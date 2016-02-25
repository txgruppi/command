package command

import "strings"

// NoHandlerFoundError is returned when no handler returns `true` for
// the `CanHandle` call.
type NoHandlerFoundError struct {
	Command interface{}
}

func (e *NoHandlerFoundError) Error() string {
	return "No handler can handle the given command"
}

// ErrorGroup is used to group errors when using the ParallelDispatcher.
//
// Since it is not possible to stop a goroutine after it is started the
// ParallelDispatcher captures all errors and group them in the ErrorGroup
// struct.
type ErrorGroup struct {
	Errors []error
}

func (e *ErrorGroup) Error() string {
	strs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		strs[i] = err.Error()
	}
	return strings.Join(strs, "\n")
}
