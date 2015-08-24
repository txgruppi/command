package command

type Handler interface {
	CanHandle(cmd interface{}) bool
	Handle(cmd interface{}, dispatcher Dispatcher) error
}
