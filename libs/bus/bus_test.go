package bus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

type Payload struct {
	Message string
}

func TestEventBus_SubscribeAndPublish(t *testing.T) {

	var id Next = uuid.NewString
	bus, err := NewBus(id)
	if err != nil {
		panic(err)
	}
	//var ctx = context.Background()

	//bus, err := NewBus(func() string { return fmt.Sprintf("%d", time.Now().UnixNano()) })
	//if err != nil {
	//	t.Fatalf("Failed to create event bus: %v", err)
	//}

	start := time.Now()

	for i := 0; i < 10000; i++ {
		topic := fmt.Sprintf("topic%d", i)
		handlerUUID := fmt.Sprintf("handler%d", i)

		bus.RegisterTopics(topic)
		handleFunc := func(ctx context.Context, event Event) {
			if payload, ok := event.Data.(*Payload); ok {
				fmt.Printf("Received on %s: %s\n", event.Topic, payload.Message)
			} else {
				fmt.Printf("Failed to parse payload on %s\n", event.Topic)
			}
		}
		handler := Handler{
			Handle:  handleFunc,
			Matcher: topic,
		}
		bus.RegisterHandler(handlerUUID, handler)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken to add 1000 subscriptions: %v\n", duration)

	// Test publishing and receiving a message
	bus.Emit(context.Background(), "topic1", &Payload{Message: "Hello, World!"})

	// Give some time for the message to be processed
	time.Sleep(1 * time.Second)
}
