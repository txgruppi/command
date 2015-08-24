package command

import "strings"

type NoHandlerFoundError struct {
	Command interface{}
}

func (e *NoHandlerFoundError) Error() string {
	return "No handler can handle the given command"
}

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
