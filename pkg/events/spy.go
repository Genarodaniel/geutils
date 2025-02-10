package events

import (
	"sync"
	"time"
)

type EventHandlerSpy struct {
	Event EventInterface
}

type EventSpy struct {
	GetNameResponse     string
	GetDateTimeResponse time.Time
	GetPayloadResponse  any
}

func (e EventSpy) GetName() string {
	return e.GetNameResponse
}

func (e EventSpy) GetDateTime() time.Time {
	return e.GetDateTimeResponse
}

func (e EventSpy) GetPayload() any {
	return e.GetPayloadResponse
}

func (eh EventHandlerSpy) Handle(event EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
}
