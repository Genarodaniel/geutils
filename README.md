# Event-Driven Worker

This project implements a generic event-driven worker that can be used with various message brokers. It provides an event-dispatching mechanism where users only need to implement the event interface and handlers. The **EventDispatcher** is already implemented, making it easy to register and dispatch events.

## Features
- Supports integration with different message brokers (Kafka, RabbitMQ, NATS, etc.)
- Dispatches events based on **message keys**
- Processes messages **concurrently** using goroutines
- Ensures partition-aware or queue-aware message handling

---

## Getting Started

### **1. Install Dependencies**
Ensure you have Go installed (1.18+ recommended) and initialize your project:

```sh
go mod init your_project
```

### **2. Define Your Event Interface**
Each event must implement the `EventInterface`:

```go
package events

import "time"

type EventInterface interface {
    GetName() string
    GetDateTime() time.Time
    GetPayload() []byte
}
```

Example event struct:

```go
type CreateOrderEvent struct {
    DateTime time.Time
    Payload  []byte
}

func (e *CreateOrderEvent) GetName() string       { return "createOrder" }
func (e *CreateOrderEvent) GetDateTime() time.Time { return e.DateTime }
func (e *CreateOrderEvent) GetPayload() []byte     { return e.Payload }
```

---

### **3. Implement an Event Handler**
Handlers must implement the `EventHandlerInterface`:

```go
package events

import (
    "fmt"
    "sync"
)

type OrderEventHandler struct {}

func (h *OrderEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Processing Order Event:", string(event.GetPayload()))
}
```

---

### **4. Run the Worker**
The worker listens for messages from the broker, determines the event type by **message key**, and dispatches them accordingly.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "your_project/events"
)

func main() {
    dispatcher := events.NewEventDispatcher()
    orderHandler := &events.OrderEventHandler{}

    ctx := context.Background()
    dispatcher.Register(ctx, "createOrder", orderHandler)

    for {
        messages := fetchMessagesFromBroker() // Implement this based on your broker
        var wg sync.WaitGroup

        for _, msg := range messages {
            wg.Add(1)
            go func(msg Message) {
                defer wg.Done()
                eventKey := msg.Key
                switch eventKey {
                case "createOrder":
                    event := &events.CreateOrderEvent{
                        DateTime: time.Now(),
                        Payload:  msg.Value,
                    }
                    fmt.Println("Received createOrder event:", string(msg.Value))
                    dispatcher.Dispatch(ctx, event)
                default:
                    fmt.Println("Unknown event key:", eventKey)
                }
            }(msg)
        }
        wg.Wait()
    }
}
```

---

## Running the Worker
1. Start your message broker (Kafka, RabbitMQ, etc.) and configure topics/queues.
2. Run the worker:

```sh
go run main.go
```

3. Produce messages to the broker using your preferred method.

---

## Extending the Worker
To add new event types:
1. **Create a new event struct** implementing `EventInterface`.
2. **Implement a handler** that satisfies `EventHandlerInterface`.
3. **Register the event and handler** in the worker.

---

## Conclusion
This worker provides a **generic event-driven architecture**, allowing integration with different message brokers. The **EventDispatcher** is already implemented, so you only need to define your events and handlers.

🚀 Happy coding!

