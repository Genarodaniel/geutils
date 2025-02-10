package events

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	getEventName := "order.create"

	t.Run("Should return error if the event are already registered", func(t *testing.T) {
		handlerSpy := EventHandlerSpy{
			Event: &EventSpy{},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		err = dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.NotNil(t, err)
		assert.EqualError(t, err, ErrHandlerAlreadyRegistered.Error())

	})

	t.Run("Should successful register a new event", func(t *testing.T) {
		handlerSpy := EventHandlerSpy{
			Event: &EventSpy{},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)
		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.True(t, created)
	})
}

func TestClear(t *testing.T) {
	gin.SetMode(gin.TestMode)
	getEventName := "order.create"

	t.Run("Should clear the event handler", func(t *testing.T) {
		handlerSpy := EventHandlerSpy{
			Event: &EventSpy{},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.True(t, created)

		dispatcher.Clear(ctx)

		created = dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.False(t, created)

	})
}

func TestHas(t *testing.T) {
	gin.SetMode(gin.TestMode)
	getEventName := "order.create"

	handlerSpy := EventHandlerSpy{
		Event: &EventSpy{},
	}

	t.Run("Should return true when the event recently created", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.True(t, created)

		dispatcher.Clear(ctx)

	})

	t.Run("Should return false when there's no event created", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.False(t, created)
	})
}

func TestRemove(t *testing.T) {
	gin.SetMode(gin.TestMode)
	getEventName := "order.create"

	handlerSpy := EventHandlerSpy{
		Event: &EventSpy{},
	}

	t.Run("Should not remove if the event is the same and the handler is different", func(t *testing.T) {
		handlerSpy2 := EventHandlerSpy{
			Event: &EventSpy{
				GetNameResponse: "order.updated",
			},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()
		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		dispatcher.Remove(ctx, getEventName, handlerSpy2)

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.True(t, created)

		dispatcher.Clear(ctx)

	})

	t.Run("Should remove if the event exists with the same handler", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()

		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		dispatcher.Remove(ctx, getEventName, handlerSpy)

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.False(t, created)

		dispatcher.Clear(ctx)

	})
}

func TestDispatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	getEventName := "order.create"
	timeResponse := time.Now()
	payloadResponse := map[string]string{
		"order_id": uuid.NewString(),
	}

	handlerSpy := EventHandlerSpy{
		Event: &EventSpy{
			GetNameResponse:     getEventName,
			GetDateTimeResponse: timeResponse,
			GetPayloadResponse:  payloadResponse,
		},
	}

	t.Run("Should return error if the event are not dispatched", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.False(t, created)

		err := dispatcher.Dispatch(ctx, handlerSpy.Event)
		assert.EqualError(t, err, ErrEventNotRegistered.Error())
	})

	t.Run("Should dispatch the event", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		dispatcher := NewEventDispatcher()

		err := dispatcher.Register(ctx, getEventName, handlerSpy)
		assert.Nil(t, err)

		created := dispatcher.Has(ctx, getEventName, handlerSpy)
		assert.True(t, created)

		err = dispatcher.Dispatch(ctx, handlerSpy.Event)
		assert.Nil(t, err)

	})

}
