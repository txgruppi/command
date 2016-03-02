package command

// NoHandlerFoundError is returned when no handler returns `true` for
// the `CanHandle` call.
type NoHandlerFoundError struct {
	Command interface{}
}

func (e *NoHandlerFoundError) Error() string {
	return "No handler can handle the given command"
}
