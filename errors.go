package command

import "fmt"

// NoHandlerFoundError is returned when no handler returns `true` for
// the `CanHandle` call.
type NoHandlerFoundError struct {
	Command interface{}
}

func (e *NoHandlerFoundError) Error() string {
	return fmt.Sprintf("No handler can handle %T", e.Command)
}
