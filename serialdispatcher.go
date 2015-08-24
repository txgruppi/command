package command

func NewSerialDispatcher(handlers []Handler) Dispatcher {
	return &SerialDispatcher{
		handlers: handlers,
	}
}

type SerialDispatcher struct {
	handlers []Handler
}

func (d *SerialDispatcher) Dispatch(cmd interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	found := false
	for _, handler := range d.handlers {
		if handler == nil {
			continue
		}
		if !handler.CanHandle(cmd) {
			continue
		}
		found = true
		if err = handler.Handle(cmd, d); err != nil {
			return
		}
	}

	if !found {
		return &NoHandlerFoundError{
			Command: cmd,
		}
	}

	return
}

func (d *SerialDispatcher) DispatchOptional(cmd interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	err = d.Dispatch(cmd)
	switch err.(type) {
	case *NoHandlerFoundError:
		return nil
	default:
		return err
	}
}
