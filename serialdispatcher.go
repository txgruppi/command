package command

import "sync"

func NewSerialDispatcher(handlers []Handler) Dispatcher {
	return &SerialDispatcher{
		handlers: handlers,
		mutex:    sync.RWMutex{},
	}
}

type SerialDispatcher struct {
	handlers []Handler
	mutex    sync.RWMutex
}

func (d *SerialDispatcher) AppendHandlers(handlers ...Handler) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

Loop:
	for _, newHandler := range handlers {
		for _, existingHandler := range d.handlers {
			if newHandler == existingHandler {
				break Loop
			}
		}
		d.handlers = append(d.handlers, newHandler)
	}
}

func (d *SerialDispatcher) Dispatch(cmd interface{}) (err error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

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
	d.mutex.RLock()
	defer d.mutex.RUnlock()

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
