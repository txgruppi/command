package command

type Dispatcher interface {
	Dispatch(cmd interface{}) error
	DispatchOptional(cmd interface{}) error
}
