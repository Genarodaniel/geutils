package events

import (
	"context"
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")
var ErrEventNotRegistered = errors.New("event not registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: map[string][]EventHandlerInterface{},
	}
}

func (ed *EventDispatcher) Dispatch(ctx context.Context, event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
		return nil
	}

	return ErrEventNotRegistered
}

func (ed *EventDispatcher) Register(ctx context.Context, event string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[event]; ok {
		for _, h := range ed.handlers[event] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	ed.handlers[event] = append(ed.handlers[event], handler)
	return nil
}

func (ed *EventDispatcher) Clear(ctx context.Context) {
	ed.handlers = map[string][]EventHandlerInterface{}
}

func (ed *EventDispatcher) Remove(ctx context.Context, event string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[event]; ok {
		for key, h := range ed.handlers[event] {
			if h == handler {
				ed.handlers[event] = append(ed.handlers[event][:key], ed.handlers[event][key+1:]...)
			}
		}
	}

	return nil
}

func (ed *EventDispatcher) Has(ctx context.Context, event string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[event]; ok {
		for _, h := range ed.handlers[event] {
			if h == handler {
				return true
			}
		}
	}

	return false
}
