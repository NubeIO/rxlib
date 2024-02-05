package rxlib

import "time"

type EventBusMessage struct {
	Data                      any           `json:"data"`
	ObjectUUID                string        `json:"objectUUID,omitempty"`
	ObjectID                  string        `json:"objectID,omitempty"`
	Topic                     string        `json:"topic,omitempty"`
	ResponseTopic             string        `json:"responseTopic,omitempty"`
	UnsubsribeOnResponseTopic bool          `json:"unsubsribe,omitempty"` // used for when we want to use the EventBus PuplishWait and we unsubsribe to the ResponseTopic
	Payload                   *Payload      `json:"payload,omitempty"`
	Timeout                   time.Duration `json:"timeout,omitempty"`
}

type EventBusCallback struct {
	EventBusCallback func(msg *EventBusMessage) `json:"-"` // used for the evntbus
}

func NewEventBusCallback(f func(msg *EventBusMessage)) *EventBusCallback {
	return &EventBusCallback{
		EventBusCallback: f,
	}
}
