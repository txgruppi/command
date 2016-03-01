package command_test

import (
	"github.com/txgruppi/command"
)

type Command uint16

const (
	CommandA Command = iota + 1
	CommandB
	CommandC
)

func NewCanHandleCallbackForCommand(cmd Command) func(interface{}) bool {
	return func(icmd interface{}) bool {
		if icmd == nil {
			return false
		}
		if c, ok := icmd.(Command); ok {
			return c == cmd
		}
		return false
	}
}

type TestHandler struct {
	CanHandleCallback  func(interface{}) bool
	HandleCallback     func(interface{}, command.Dispatcher) error
	CanHandleCallCount int
	HandleCallCount    int
}

func (h *TestHandler) CanHandle(cmd interface{}) bool {
	h.CanHandleCallCount++
	if h.CanHandleCallback == nil {
		return false
	}
	return h.CanHandleCallback(cmd)
}

func (h *TestHandler) Handle(cmd interface{}, dispatcher command.Dispatcher) error {
	h.HandleCallCount++
	if h.HandleCallback == nil {
		return nil
	}
	return h.HandleCallback(cmd, dispatcher)
}
