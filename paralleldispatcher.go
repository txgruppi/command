package command

import (
	"sync"
	"sync/atomic"
)

// NewParallelDispatcher creates a new PrallelDispatcher with the given handlers
func NewParallelDispatcher(handlers []Handler) Dispatcher {
	return &ParallelDispatcher{
		handlers: handlers,
		mutex:    sync.RWMutex{},
	}
}

// ParallelDispatcher is a command dispatcher wich will run all handlers in
// parallel and wait all handlers to finish before returning.
//
// All errors returned by the handlers will be grouped in a `ErrorGroup`
// instance.
//
// This dispatcher is *thread safe*.
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

	var found int32
	wg := &sync.WaitGroup{}
	errCh := make(chan error, len(d.handlers))
	for _, handler := range d.handlers {
		wg.Add(1)
		go d.dispatch(wg, errCh, &found, handler, cmd)
	}

	wg.Wait()
	close(errCh)

	if found != 1 {
		return &NoHandlerFoundError{
			Command: cmd,
		}
	}

	errs := []error{}
	for {
		e, ok := <-errCh
		if !ok {
			break
		}
		if e == nil {
			continue
		}
		errs = append(errs, e)
	}

	if len(errs) == 0 {
		return
	}

	err = &ErrorGroup{
		Errors: errs,
	}

	return
}

func (d *ParallelDispatcher) dispatch(wg *sync.WaitGroup, errCh chan error, found *int32, handler Handler, cmd interface{}) {
	var err error

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		errCh <- err
		wg.Done()
	}()

	if !handler.CanHandle(cmd) {
		return
	}

	atomic.StoreInt32(found, 1)

	err = handler.Handle(cmd, d)
}

func (d *ParallelDispatcher) DispatchOptional(cmd interface{}) (err error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	err = d.Dispatch(cmd)
	switch err.(type) {
	case *NoHandlerFoundError:
		return nil
	default:
		return err
	}
}
