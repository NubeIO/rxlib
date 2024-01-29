package rxlib

import (
	"fmt"
	"github.com/gookit/event"
	"sync"
)

type EventBus struct {
	manager         *event.Manager
	handlerRegistry map[string]map[string]event.ListenerFunc
	activeHandlers  map[string]bool
	mu              sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		manager:         event.NewManager("rx-lib"),
		handlerRegistry: make(map[string]map[string]event.ListenerFunc),
		activeHandlers:  make(map[string]bool),
	}
}

func (eb *EventBus) Subscribe(topic string, handlerID string, handler func(e event.Event) error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	wrappedHandler := func(e event.Event) error {
		// Check if handler is active before processing
		if _, active := eb.activeHandlers[handlerID]; active {
			return handler(e)
		}
		return nil
	}

	listenerFunc := event.ListenerFunc(wrappedHandler)
	eb.manager.On(topic, listenerFunc, event.Normal)

	if _, ok := eb.handlerRegistry[topic]; !ok {
		eb.handlerRegistry[topic] = make(map[string]event.ListenerFunc)
	}
	eb.handlerRegistry[topic][handlerID] = listenerFunc
	eb.activeHandlers[handlerID] = true
}

func (eb *EventBus) Unsubscribe(topic string, handlerID string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if listeners, ok := eb.handlerRegistry[topic]; ok {
		if _, ok := listeners[handlerID]; ok {
			fmt.Println("Deactivate handler topic:", topic)
			delete(listeners, handlerID)
			delete(eb.activeHandlers, handlerID) // Deactivate handler
		}

		if len(listeners) == 0 {
			fmt.Printf("remove listeners, current curent: %d topic %s \n", len(listeners), topic)
			delete(eb.handlerRegistry, topic)
			eb.manager.RemoveEvent(topic)
		} else {
			fmt.Printf("Cant remove listeners, current curent: %d topic %s \n", len(listeners), topic)
		}
	}
}

func (eb *EventBus) Publish(topic string, data *Payload) error {
	eventData := event.M{
		"message": data,
	}
	err, _ := eb.manager.Fire(topic, eventData)
	return err
}

func (eb *EventBus) ListSubscribers() map[string][]string {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := make(map[string][]string)
	for topic, handlers := range eb.handlerRegistry {
		for handlerID := range handlers {
			subscribers[topic] = append(subscribers[topic], handlerID)
		}
	}
	return subscribers
}

func (eb *EventBus) ListSubscribersForTopic(topic string) []string {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var subscribers []string
	if handlers, ok := eb.handlerRegistry[topic]; ok {
		for handlerID := range handlers {
			subscribers = append(subscribers, handlerID)
		}
	}
	return subscribers
}

func (eb *EventBus) GetSubscriberCount() int {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	count := 0
	for _, handlers := range eb.handlerRegistry {
		count += len(handlers)
	}
	return count
}
