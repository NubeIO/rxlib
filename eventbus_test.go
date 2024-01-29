package rxlib

import (
	"fmt"
	"github.com/gookit/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//type DummyMessage struct {
//	Data string
//}

func TestNewEventBus(t *testing.T) {
	eb := NewEventBus()
	assert.NotNil(t, eb)

	// Define a topic
	topic := "testTopic"

	// Define a channel to receive message
	received := make(chan *Message, 1)

	// Subscribe to the topic
	eb.Subscribe(topic, "testHandler", func(e event.Event) error {
		msg, ok := e.Get("message").(*Message)
		if ok {
			received <- msg
		}
		return nil
	})

	// Publish a message to the topic
	testMessage := &Message{ObjectUUID: "Hello, EventBus!"}
	eb.Publish(topic, testMessage)

	// Wait and check if we received the message
	select {
	case msg := <-received:
		fmt.Println(msg, 999999999)
		assert.Equal(t, testMessage, msg)
	case <-time.After(1 * time.Second):
		t.Errorf("Did not receive message on topic %s", topic)
	}

}
