package rxlib

type EventBusMessage struct {
	Data       any    `json:"data"`
	ObjectUUID string `json:"objectUUID,omitempty"`
	ObjectID   string `json:"objectID,omitempty"`
	Topic      string `json:"topic,omitempty"`
}

type EventBusCallback struct {
	EventBusCallback func(msg *EventBusMessage) `json:"-"` // used for the evntbus
}

func NewEventBusCallback(f func(msg *EventBusMessage)) *EventBusCallback {
	return &EventBusCallback{
		EventBusCallback: f,
	}
}
