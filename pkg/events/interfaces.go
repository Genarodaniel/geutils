package events

import (
	"context"
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(ctx context.Context, event string, handler EventHandlerInterface) error
	Dispatch(ctx context.Context, event EventInterface) error
	Remove(ctx context.Context, event string, handler EventHandlerInterface) error
	Has(ctx context.Context, event string, handler EventHandlerInterface) bool
	Clear(ctx context.Context) error
}
