package event

import (
	"fmt"
	"github.com/NubeIO/rxlib/payload"
	"sync"
)

type EventBus interface {
	Subscribe(topic string, callback func(topic string, data *payload.Payload, err error))
	Publish(topic string, data *payload.Payload) error
	Unsubscribe(topic string, callback func(topic string, data *payload.Payload, err error))
}

type eventBus struct {
	subscribers map[string][]func(topic string, data *payload.Payload, err error)
	pubChannel  chan *eventMessage
	lock        sync.RWMutex
}

type eventMessage struct {
	topic string
	data  *payload.Payload
}

func NewEventBus() EventBus {
	bus := &eventBus{
		subscribers: make(map[string][]func(topic string, data *payload.Payload, err error)),
		pubChannel:  make(chan *eventMessage, 100),
	}
	go bus.start()
	return bus
}

func (b *eventBus) Subscribe(topic string, callback func(topic string, data *payload.Payload, err error)) {
	b.lock.Lock()
	defer b.lock.Unlock()
	fmt.Printf("EVENTBUS Subscribe topic: %s \n", topic)
	b.subscribers[topic] = append(b.subscribers[topic], callback)
	fmt.Printf("Current Subscribers for topic %s: %v\n", topic, b.subscribers[topic])
}

func (b *eventBus) Publish(topic string, data *payload.Payload) error {
	fmt.Printf("EVENTBUS Publish topic: %s\n", topic)
	b.pubChannel <- &eventMessage{
		topic: topic,
		data:  data,
	}
	return nil
}

func (b *eventBus) Unsubscribe(topic string, callback func(topic string, data *payload.Payload, err error)) {
	b.lock.Lock()
	defer b.lock.Unlock()
	fmt.Printf("EVENTBUS Unsubscribe topic: %s \n", topic)
	// Find and remove the specific callback
	if subscribers, ok := b.subscribers[topic]; ok {
		for i, sub := range subscribers {
			if fmt.Sprintf("%p", sub) == fmt.Sprintf("%p", callback) {
				b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		// If there are no more subscribers for the topic, delete the entry
		if len(b.subscribers[topic]) == 0 {
			delete(b.subscribers, topic)
		}
	}
	fmt.Printf("Subscribers after Unsubscribe for topic %s: %v\n", topic, b.subscribers[topic])
}

func (b *eventBus) start() {
	for msg := range b.pubChannel {
		b.lock.RLock()
		if subscribers, ok := b.subscribers[msg.topic]; ok {
			for _, subscriber := range subscribers {
				go subscriber(msg.topic, msg.data, nil)
			}
		}
		b.lock.RUnlock()
	}
}
