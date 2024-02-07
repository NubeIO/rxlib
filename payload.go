package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/priority"
	"time"
)

func NewPayload() *Payload {
	return &Payload{}
}

type ObjectInvoke struct {
	FromObjectUUID string `json:"fromObjectUUID"`
	Action         string `json:"action"`
	Data           any    `json:"data"`
}

type ObjectInvokeResponse struct {
	Data any `json:"data"`
}

type Payload struct {
	PortPayload *PortPayload `json:"portPayload,omitempty"`
	// used for mapping
	Mapping *Mapping `json:"mapping,omitempty"`
	// generic eventbus message
	EventBusPayload *EventBusPayload `json:"eventBusPayload,omitempty"`
}

type PortPayload struct {
	Port       *Port       `json:"port,omitempty"`
	Connection *Connection `json:"connection,omitempty"`
	ObjectUUID string      `json:"objectUUID,omitempty"`
	ObjectID   string      `json:"objectID,omitempty"`
}

type EventBusPayload struct {
	ObjectUUID                 string        `json:"objectUUID,omitempty"`
	ObjectID                   string        `json:"objectID,omitempty"`
	Topic                      string        `json:"topic,omitempty"`
	ResponseTopic              string        `json:"responseTopic,omitempty"`
	UnsubscribeOnResponseTopic bool          `json:"unsubscribe,omitempty"` // used for when we want to use the EventBus PublishWait and we unsubscribe to the ResponseTopic
	Timeout                    time.Duration `json:"timeout,omitempty"`
	ExpectedData               string        `json:"expectedData,omitempty"` // make it easy for an object to decode in incoming data; eg string, map[], user
	Payload                    any           `json:"payload,omitempty"`
}

type Mapping struct {
	ManagerUUID       string            `json:"managerUUID,omitempty"`
	NetworkUUID       string            `json:"networkUUID,omitempty"`
	MapperUUID        string            `json:"mapperUUID,omitempty"`
	Data              any               `json:"data,omitempty"`
	ExpectedData      string            `json:"expectedData,omitempty"` // make it easy for an object to decode in incoming data; eg string, map[], user
	PrimitivesPayload PrimitivesPayload `json:"primitivesPayload,omitempty"`
}

type PrimitivesPayload struct {
	DataType priority.Type      `json:"dataType,omitempty"`
	Priority *priority.Priority `json:"priority,omitempty"`
	Symbol   *string            `json:"symbol,omitempty"`
}

type PayloadValue struct {
	Value     any `json:"value"`
	Timestamp *time.Time
}

func (p *Payload) GetPayload() *Payload {
	return p
}

// ----------------PORT------------------

func (p *Payload) NewPortPayload() *Payload {
	p.PortPayload = &PortPayload{}
	return p
}

func (p *Payload) IsPortNil() bool {
	if p.PortPayload == nil {
		return true
	}
	if p.PortPayload.Port == nil {
		return true
	}
	return false
}

func (p *Payload) GetPort() *Port {
	if p.IsPortNil() {
		return nil
	}
	return p.PortPayload.Port
}

func (p *Payload) SetPort(port *Port) *Payload {
	if p.IsPortNil() {
		p.PortPayload = &PortPayload{}
	}
	p.PortPayload.Port = port
	return p
}

func (p *Payload) SetPortObjectDetails(objectID, objectUUID string) *Payload {
	if p.IsPortNil() {
		p.PortPayload = &PortPayload{}
	}
	p.PortPayload.ObjectID = objectID
	p.PortPayload.ObjectUUID = objectUUID
	return p
}

func (p *Payload) GetValue() (any, error) {
	if p.IsPortNil() {
		return nil, fmt.Errorf("port is empty")
	}
	return p.GetPort().Data, nil
}

func (p *Payload) IsPriorityNil() bool {
	if p.IsPortNil() {
		return true
	}
	if p.PortPayload.Port.DataPriority == nil {
		return true
	}
	return false
}

func (p *Payload) GetDataPriority() (*priority.Priority, error) {
	if p.IsPriorityNil() {
		return nil, fmt.Errorf("port has no priority set")
	}
	return p.GetPort().DataPriority.Priority, nil
}

func (p *Payload) GetPortName() string {
	return p.GetPort().Name
}

// ----------------CONNECTION------------------

func (p *Payload) IsConnectionNil() bool {
	if p.PortPayload == nil {
		return true
	}
	if p.PortPayload.Connection == nil {
		return true
	}
	return false
}

func (p *Payload) GetConnection() *Connection {
	if p.IsConnectionNil() {
		return nil
	}
	return p.PortPayload.Connection
}

// ----------------EVENTBUS------------------

func (p *Payload) SetEventbusPayload(body *EventBusPayload) *Payload {
	p.EventBusPayload = body
	return p
}

func (p *Payload) NewEventbusPayload(topic string, payload any) *Payload {
	p.EventBusPayload = &EventBusPayload{
		Topic:   topic,
		Payload: payload,
	}
	return p
}

func (p *Payload) IsEventBusPayloadNil() bool {
	if p.EventBusPayload == nil {
		return true
	}
	return false
}

func (p *Payload) GetEventbusObjectID() string {
	if p.IsEventBusPayloadNil() {
		return ""
	}
	return p.EventBusPayload.ObjectID
}

func (p *Payload) GetEventbusObjectUUID() string {
	if p.IsEventBusPayloadNil() {
		return ""
	}
	return p.EventBusPayload.ObjectUUID
}

func (p *Payload) SetEventbusObjectUUID(value string) {
	if p.IsEventBusPayloadNil() {
		return
	}
	p.EventBusPayload.ObjectUUID = value
}

func (p *Payload) GetExpectedData() string {
	if p.IsEventBusPayloadNil() {
		return ""
	}
	return p.EventBusPayload.ExpectedData
}

func (p *Payload) GetEventBusPayload() *EventBusPayload {
	if p.IsEventBusPayloadNil() {
		return nil
	}
	return p.EventBusPayload
}

func (p *Payload) UnsubscribeOnResponseTopic() bool {
	if p.IsEventBusPayloadNil() {
		return false
	}
	return p.GetEventBusPayload().UnsubscribeOnResponseTopic
}

// ----------------EVENTBUS------------------

func (p *Payload) NewMapping() *Payload {
	p.Mapping = &Mapping{}
	return p
}

func (p *Payload) GetMapping() *Mapping {
	if p.IsMappingNil() {
		return nil
	}
	return p.Mapping
}

func (p *Payload) SetMappingDetails(managerUUID, networkUUID, mapperUUID string) *Payload {
	if p.Mapping == nil {
		p.Mapping = &Mapping{
			ManagerUUID: managerUUID,
			NetworkUUID: networkUUID,
			MapperUUID:  mapperUUID,
		}
	} else {
		p.Mapping.ManagerUUID = managerUUID
		p.Mapping.NetworkUUID = networkUUID
		p.Mapping.MapperUUID = mapperUUID
	}
	return p
}

func (p *Payload) SetMappingData(expectedData string, data any) *Payload {
	p.Mapping = &Mapping{
		ExpectedData: expectedData,
		Data:         data,
	}
	return p
}

func (p *Payload) IsMappingNil() bool {
	if p.Mapping == nil {
		return true
	}
	return false
}

func (p *Payload) GetMappingManagerUUID() string {
	if p.IsMappingNil() {
		return ""
	}
	return p.GetMapping().ManagerUUID
}

func (p *Payload) GetMappingMapperUUID() string {
	if p.IsMappingNil() {
		return ""
	}
	return p.GetMapping().MapperUUID
}

func (p *Payload) GetMappingNetworkUUID() string {
	if p.IsMappingNil() {
		return ""
	}
	return p.GetMapping().NetworkUUID
}

func (p *Payload) GetMappingData() any {
	if p.IsMappingNil() {
		return ""
	}
	return p.GetMapping().Data
}
