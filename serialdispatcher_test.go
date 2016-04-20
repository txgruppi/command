package command_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/txgruppi/command"
)

func TestSerialDispatcher(t *testing.T) {
	Convey("SerialDispatcher", t, func() {
		handlerA := &TestHandler{
			CanHandleCallback: NewCanHandleCallbackForCommand(CommandA),
		}
		handlerB := &TestHandler{
			CanHandleCallback: NewCanHandleCallbackForCommand(CommandB),
		}
		handlers := []command.Handler{handlerA, handlerB}
		dispatcher := command.NewSerialDispatcher(handlers)

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
				So(err.Error(), ShouldEqual, "No handler can handle command_test.Command")
			})

			Convey("it should recover from panic", func() {
				canHandleErr := errors.New("CanHandle panic")
				handleErr := errors.New("Handle panic")
				handlerA.CanHandleCallback = func(interface{}) bool { panic(canHandleErr) }
				err := dispatcher.Dispatch(CommandA)
				So(err.Error(), ShouldEqual, canHandleErr.Error())
				handlerA.CanHandleCallback = nil
				handlerB.HandleCallback = func(interface{}, command.Dispatcher) error { panic(handleErr) }
				err = dispatcher.Dispatch(CommandB)
				So(err.Error(), ShouldEqual, handleErr.Error())
			})

			Convey("it should wait a handler to finish before continue (serial dispatch)", func() {
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
					return true
				}
				err := dispatcher.Dispatch(CommandB)
				So(err, ShouldBeNil)
				So(callOrder, ShouldResemble, []int{1, 2})
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
