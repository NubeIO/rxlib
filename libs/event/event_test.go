package event

import (
	"fmt"
	"github.com/NubeIO/rxlib/payload"
	"testing"
	"time"
)

func TestNewEventBus(t *testing.T) {
	bus := NewEventBus()

	bus.Subscribe("topic1", func(topic string, data *payload.Payload, err error) {
		fmt.Println("Received on 11111:", data, topic)
	})

	bus.Publish("topic1", &payload.Payload{
		FromPortID:     "111",
		FromObjectUUID: "",
		PortValue:      nil,
	})

	// Unsubscribe from topic1

	// This message will not be received by any subscriber
	//bus.Publish("topic1", &Data{Any: "Hello 2, World!"})

	time.Sleep(1 * time.Second)

	bus.Unsubscribe("topic1")

}

func TestNewEventBus1000(t *testing.T) {
	bus := NewEventBus()

	// Start timing
	start := time.Now()

	// Add 1000 subscriptions
	for i := 0; i < 10000000; i++ {
		bus.Subscribe(fmt.Sprintf("topic%d", i), func(topic string, data *payload.Payload, err error) {
			fmt.Printf("Received on topic%d: %v\n", i, data)
		})
	}

	// Stop timing and print the duration
	duration := time.Since(start)
	fmt.Printf("Time taken to add 1000 subscriptions: %v\n", duration)

	// Test publishing and receiving a message
	bus.Publish("topic1", &payload.Payload{})

	// Give some time for the message to be processed
	time.Sleep(1 * time.Second)

	// Unsubscribe from topic1
	bus.Unsubscribe("topic1")

}
