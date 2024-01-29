package rxlib

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
)

type Message struct {
	Port       *Port  `json:"port"`
	ObjectUUID string `json:"objectUUID"`
	ObjectID   string `json:"objectID"`
}

type ClosedChannel struct {
	Closed       bool
	ObjectUUID   string
	ObjectPortID string
}

// EventBus manages event subscriptions and publishes events.
type EventBus struct {
	mu             sync.Mutex
	handlers       map[string][]chan *Message
	subscribers    map[chan *Message]string
	closedChannels map[chan *Message]*ClosedChannel
	globalChan     chan *Message // Global channel for publishing and subscribing
	//WS             *WSHub
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	//ws := NewWSHub()
	//go ws.Run()
	return &EventBus{
		handlers:       make(map[string][]chan *Message),
		subscribers:    make(map[chan *Message]string), // inputs
		closedChannels: make(map[chan *Message]*ClosedChannel),
		globalChan:     make(chan *Message),
		//WS:             ws,
	}
}

func (eb *EventBus) Subscribe(topic string, ch chan *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	fmt.Println("EVENT BUS Subscribe", topic)
	eb.handlers[topic] = append(eb.handlers[topic], ch)
	eb.subscribers[ch] = topic
}

// Unsubscribe unsubscribes a channel from a topic.
func (eb *EventBus) Unsubscribe(topic string, ch chan *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	subscribers := eb.handlers[topic]
	fmt.Printf("Unsubscribed from topic: %s subscribers count: %d \n", topic, len(subscribers))
	for i, sub := range subscribers {
		if sub == ch {
			close(sub)                 // Close the channel to stop the goroutine
			subscribers[i] = nil       // Set the channel to nil
			delete(eb.subscribers, ch) // Remove the subscriber entry
			fmt.Printf("Unsubscribed delete topic: %s\n", topic)
			break
		}
	}

	eb.handlers[topic] = subscribers // Update the subscribers list for the topic

}

// MarkAsClosed was added to close the chan for an input
func (eb *EventBus) MarkAsClosed(ch chan *Message, details *ClosedChannel) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if details == nil {
		return
	}
	details.Closed = true
	eb.closedChannels[ch] = details
}

func (eb *EventBus) IsClosed(ch chan *Message) bool {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	closed, exists := eb.closedChannels[ch]
	if closed != nil {
		return false
	}
	return exists && closed.Closed
}

// Publish publishes an event to all subscribers of a topic. t
func (eb *EventBus) Publish(topic string, data *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, ch := range eb.handlers[topic] {
		go func(ch chan *Message) {
			//fmt.Printf("Publishing message to channel: %v\n", ch)
			//fmt.Printf("Message details: ObjectUUID=%s, ObjectID=%s, PortID=%s\n", data.ObjectUUID, data.ObjectID, data.Port.ID)
			// Here you could also add code to log the message details to a file or database
			if closedDetails, isClosed := eb.closedChannels[ch]; !isClosed {
				ch <- data
			} else {
				fmt.Printf("topic- is closed %v \n", closedDetails)
			}

		}(ch)
	}
}

func (eb *EventBus) DoesPublisherExist(topic string) bool {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	_, exists := eb.handlers[topic]
	return exists
}

func (eb *EventBus) AddSubscriptionExistingToPublisher(sourceUUID, sourcePortID string, subscriber chan *Message) (chan *Message, error) {
	topic := fmt.Sprintf("%s-%s", sourceUUID, sourcePortID)
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if subscriber == nil {
		subscriber = make(chan *Message)
	}

	if _, ok := eb.handlers[topic]; ok {
		// If the topic exists, add the new channel to the list of handlers for that topic
		eb.handlers[topic] = append(eb.handlers[topic], subscriber)
		eb.subscribers[subscriber] = topic
		fmt.Printf("Added new subscriber to existing topic: %s\n", topic)
	} else {
		return nil, fmt.Errorf("Topic does not exist: %s. Creating new topic.\n", topic)
	}

	return subscriber, nil
}

// GlobalSubscriber subscribes to the global channel.
func (eb *EventBus) GlobalSubscriber() chan *Message {
	return eb.globalChan
}

// GlobalPublisher publishes a message to the global channel.
func (eb *EventBus) GlobalPublisher(message *Message) {
	eb.globalChan <- message
}

// example on global GlobalPublisher
/*
	go func() {
		for i := 0; i < 11; i++ {
			message := &rxlib.Message{
				// Populate your message here
				ObjectUUID: fmt.Sprintf("Object-%d", i),
				// ... other fields ...
			}
			inst.GlobalPublisher(message)
			time.Sleep(1 * time.Second)
		}
	}()
*/

// example on global GlobalSubscriber and AddSubscriptionExistingToPublisher
/*
	newSubscriber, err := inst.AddSubscriptionExistingToPublisher("trigger-abc-output")
	fmt.Println(err)
	go func() {
		for msg := range newSubscriber {
			fmt.Printf("New Subscriber received: port: %s  value: %v object-id: %s uuid: %s \n", msg.Port.ID, msg.Port.Value, msg.ObjectID, msg.ObjectUUID)
		}
	}()

	go func() {
		globalSub := inst.GlobalSubscriber()
		for msg := range globalSub {
			fmt.Printf("Received message on global channel: %v\n", msg)
		}
	}()
*/

type EventbusTopicStats struct {
	SourceUUID   string `json:"sourceUUID,omitempty"`
	SourcePortID string `json:"sourcePortID,omitempty"`
	Topic        string `json:"topic,omitempty"`
}

// ListPublishers returns a slice of EventbusTopicStats, each representing a topic and its publisher count.
func (eb *EventBus) ListPublishers() []EventbusTopicStats {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var stats []EventbusTopicStats
	for topic, _ := range eb.handlers {
		objectUUID, portID := topicSplit(topic)
		stats = append(stats, EventbusTopicStats{
			SourceUUID:   objectUUID,
			SourcePortID: portID,
			Topic:        topic,
		})
	}
	return stats
}

// ListSubscribers returns a slice of SubscriberStats, each representing a subscriber channel and its topic.
func (eb *EventBus) ListSubscribers() []EventbusTopicStats {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	var stats []EventbusTopicStats
	for _, topic := range eb.subscribers {
		objectUUID, portID := topicSplit(topic)
		stats = append(stats, EventbusTopicStats{
			SourceUUID:   objectUUID,
			SourcePortID: portID,
			Topic:        topic,
		})
	}
	return nil
}

type WSConnection struct {
	Conn *websocket.Conn
	Send chan *Message // Channel to Send messages

}

type WSHub struct {
	Connections map[*WSConnection]bool
	Register    chan *WSConnection
	Unregister  chan *WSConnection
	Broadcast   chan *Message
	mu          sync.Mutex // Mutex to protect concurrent access to connections
}

func NewWSConnection(conn *websocket.Conn) *WSConnection {
	return &WSConnection{
		Conn: conn,
		Send: make(chan *Message),
	}
}

func topicSplit(input string) (uuid string, outputID string) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func (h *WSHub) Unsubscribe(conn *WSConnection) {
	h.mu.Lock() // Use a mutex to handle concurrent access
	defer h.mu.Unlock()

	if _, ok := h.Connections[conn]; ok {
		delete(h.Connections, conn)
		close(conn.Send)
	}
}

func NewWSHub() *WSHub {
	return &WSHub{
		Connections: make(map[*WSConnection]bool),
		Register:    make(chan *WSConnection),
		Unregister:  make(chan *WSConnection),
		Broadcast:   make(chan *Message),
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.mu.Lock()
			h.Connections[conn] = true
			h.mu.Unlock()
		case msg := <-h.Broadcast:
			h.mu.Lock()
			for conn := range h.Connections {
				select {
				case conn.Send <- msg:
				default:
					h.Unsubscribe(conn)
				}
			}
			h.mu.Unlock()
		}
	}
}
