package command

// Handler defines the interface for a command handler.
type Handler interface {
	// CanHandle should return `true` whenever the given command can be handled
	// by this Handler, otherwise it should return `false`
	CanHandle(cmd interface{}) bool

	// Handle does all the *work* realted to the given command.
	//
	// It will only be called if CanHandle returns `true` for the given command.
	Handle(cmd interface{}, dispatcher Dispatcher) error
}
