package command

import "sync"

func NewParallelDispatcher(handlers []Handler) Dispatcher {
	return &ParallelDispatcher{
		handlers: handlers,
		mutex:    sync.RWMutex{},
	}
}

type ParallelDispatcher struct {
	handlers []Handler
	mutex    sync.RWMutex
}

func (d *ParallelDispatcher) AppendHandlers(handlers ...Handler) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

Loop:
	for _, newHandler := range handlers {
		for _, existingHandler := range d.handlers {
			if newHandler == existingHandler {
				continue Loop
			}
		}
		d.handlers = append(d.handlers, newHandler)
	}
}

func (d *ParallelDispatcher) Dispatch(cmd interface{}) (err error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, len(d.handlers))
	found := false
	for _, handler := range d.handlers {
		if !handler.CanHandle(cmd) {
			continue
		}

		found = true
		wg.Add(1)
		go d.dispatch(&wg, errCh, handler, cmd)
	}

	if !found {
		return &NoHandlerFoundError{
			Command: cmd,
		}
	}

	wg.Wait()
	close(errCh)

	errs := []error{}
	for {
		e, ok := <-errCh
		if !ok {
			break
		}
		errs = append(errs, e)
	}

	err = &ErrorGroup{
		Errors: errs,
	}

	return
}

func (d *ParallelDispatcher) dispatch(wg *sync.WaitGroup, errCh chan error, handler Handler, cmd interface{}) {
	var err error

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		errCh <- err
		wg.Done()
	}()

	err = handler.Handle(cmd, d)
}

func (d *ParallelDispatcher) DispatchOptional(cmd interface{}) (err error) {
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
