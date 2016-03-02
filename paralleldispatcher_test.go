package command_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/nproc/errorgroup-go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/txgruppi/command"
)

func TestParallelDispatcher(t *testing.T) {
	Convey("ParallelDispatcher", t, func() {
		handlerA := &TestHandler{
			CanHandleCallback: NewCanHandleCallbackForCommand(CommandA),
		}
		handlerB := &TestHandler{
			CanHandleCallback: NewCanHandleCallbackForCommand(CommandB),
		}
		handlers := []command.Handler{handlerA, handlerB}
		dispatcher := command.NewParallelDispatcher(handlers)

		Convey("New", func() {
			Convey("it should create a dispatcher with the given handlers", func() {
				err := dispatcher.Dispatch(CommandA)
				So(err, ShouldBeNil)
				So(handlerA.CanHandleCallCount, ShouldEqual, 1)
				So(handlerB.CanHandleCallCount, ShouldEqual, 1)
				So(handlerA.HandleCallCount, ShouldEqual, 1)
				So(handlerB.HandleCallCount, ShouldEqual, 0)

				err = dispatcher.Dispatch(CommandB)
				So(err, ShouldBeNil)
				So(handlerA.CanHandleCallCount, ShouldEqual, 2)
				So(handlerB.CanHandleCallCount, ShouldEqual, 2)
				So(handlerA.HandleCallCount, ShouldEqual, 1)
				So(handlerB.HandleCallCount, ShouldEqual, 1)
			})
		})

		Convey("AppendHandlers", func() {
			Convey("it should append the given handler to the existing handlers", func() {
				handlerC := &TestHandler{
					CanHandleCallback: NewCanHandleCallbackForCommand(CommandC),
				}
				dispatcher.AppendHandlers(handlerC)

				err := dispatcher.Dispatch(CommandA)
				So(err, ShouldBeNil)
				So(handlerA.CanHandleCallCount, ShouldEqual, 1)
				So(handlerB.CanHandleCallCount, ShouldEqual, 1)
				So(handlerC.CanHandleCallCount, ShouldEqual, 1)
				So(handlerA.HandleCallCount, ShouldEqual, 1)
				So(handlerB.HandleCallCount, ShouldEqual, 0)
				So(handlerC.HandleCallCount, ShouldEqual, 0)

				err = dispatcher.Dispatch(CommandB)
				So(err, ShouldBeNil)
				So(handlerA.CanHandleCallCount, ShouldEqual, 2)
				So(handlerB.CanHandleCallCount, ShouldEqual, 2)
				So(handlerC.CanHandleCallCount, ShouldEqual, 2)
				So(handlerA.HandleCallCount, ShouldEqual, 1)
				So(handlerB.HandleCallCount, ShouldEqual, 1)
				So(handlerC.HandleCallCount, ShouldEqual, 0)

				err = dispatcher.Dispatch(CommandC)
				So(err, ShouldBeNil)
				So(handlerA.CanHandleCallCount, ShouldEqual, 3)
				So(handlerB.CanHandleCallCount, ShouldEqual, 3)
				So(handlerC.CanHandleCallCount, ShouldEqual, 3)
				So(handlerA.HandleCallCount, ShouldEqual, 1)
				So(handlerB.HandleCallCount, ShouldEqual, 1)
				So(handlerC.HandleCallCount, ShouldEqual, 1)
			})
		})

		Convey("Dispatch", func() {
			Convey("it should dispatch the command and return the handler's return value", func() {
				expected := errors.New("This is the expected error")
				handlerA.HandleCallback = func(interface{}, command.Dispatcher) error { return expected }
				err := dispatcher.Dispatch(CommandA)
				So(err.Error(), ShouldEqual, expected.Error())
			})

			Convey("it should dispatch the command and return the handler's return value only if the command can be handled", func() {
				notExpected := errors.New("This error should not be returned")
				handlerA.HandleCallback = func(interface{}, command.Dispatcher) error { return notExpected }
				err := dispatcher.Dispatch(CommandB)
				So(err, ShouldBeNil)
			})

			Convey("it should return an error if no handler can handle the given command", func() {
				err := dispatcher.Dispatch(CommandC)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "No handler can handle the given command")
			})

			Convey("it should recover from panic", func() {
				canHandleErr := errors.New("CanHandle panic")
				handleErr := errors.New("Handle panic")
				handlerA.CanHandleCallback = func(interface{}) bool { panic(canHandleErr) }
				err := dispatcher.Dispatch(CommandB)
				So(err.Error(), ShouldEqual, canHandleErr.Error())
				handlerA.CanHandleCallback = nil
				handlerB.HandleCallback = func(interface{}, command.Dispatcher) error { panic(handleErr) }
				err = dispatcher.Dispatch(CommandB)
				So(err.Error(), ShouldEqual, handleErr.Error())
			})

			Convey("it should wait a handler to finish before continue (parallel dispatch)", func() {
				callOrder := []int{}
				lock := sync.Mutex{}
				handlerA.CanHandleCallback = func(interface{}) bool {
					time.Sleep(10 * time.Millisecond)
					lock.Lock()
					defer lock.Unlock()
					callOrder = append(callOrder, 1)
					return false
				}
				handlerB.CanHandleCallback = func(interface{}) bool {
					lock.Lock()
					defer lock.Unlock()
					callOrder = append(callOrder, 2)
					return false
				}
				dispatcher.Dispatch(CommandB)
				So(callOrder, ShouldResemble, []int{2, 1})
			})

			Convey("it should group errors", func() {
				errA := errors.New("Error A")
				errB := errors.New("Error B")
				handlerA.HandleCallback = func(interface{}, command.Dispatcher) error { panic(errA) }
				handlerB.HandleCallback = func(interface{}, command.Dispatcher) error { return errB }
				handlerB.CanHandleCallback = NewCanHandleCallbackForCommand(CommandA)
				err := dispatcher.Dispatch(CommandA)
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, &errorgroup.ErrorGroup{})
				errGroup := err.(*errorgroup.ErrorGroup)
				So(len(errGroup.Errors), ShouldEqual, 2)
				So(errGroup.Errors, ShouldContain, errA)
				So(errGroup.Errors, ShouldContain, errB)
			})
		})

		Convey("DispatchOptional", func() {
			Convey("it should dispatch the command and return the handler's return value", func() {
				expected := errors.New("This is the expected error")
				handlerA.HandleCallback = func(interface{}, command.Dispatcher) error { return expected }
				err := dispatcher.DispatchOptional(CommandA)
				So(err.Error(), ShouldEqual, expected.Error())
			})

			Convey("it should return nil if no handler can handle the given command", func() {
				err := dispatcher.DispatchOptional(CommandC)
				So(err, ShouldBeNil)
			})
		})
	})
}
