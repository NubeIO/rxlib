package event

import (
	"github.com/NubeIO/rxlib/payload"
	"sync"
)

type EventBus interface {
	Subscribe(topic string, callback func(topic string, data *payload.Payload, err error))
	Publish(topic string, data *payload.Payload) error
	Unsubscribe(topic string)
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
	b.subscribers[topic] = append(b.subscribers[topic], callback)
}

func (b *eventBus) Publish(topic string, data *payload.Payload) error {
	b.pubChannel <- &eventMessage{
		topic: topic,
		data:  data,
	}
	return nil
}

func (b *eventBus) Unsubscribe(topic string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	// Delete all subscribers for the topic
	delete(b.subscribers, topic)
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
