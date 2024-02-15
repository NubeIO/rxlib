package rxlib

import (
	"encoding/json"
	"github.com/NubeIO/rxlib/priority"
	"strings"
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

type ObjectCommandResponse struct {
	Data any `json:"data"`
}

type Payload struct {
	DataPayload *DataPayload `json:"data,omitempty"`
	// used for mapping
	Mapping *Mapping `json:"mapping,omitempty"`
	// generic eventbus message
	EventBusPayload *EventBusPayload `json:"eventBusPayload,omitempty"`
}

type Command struct {
	NameSpace  string `json:"nameSpace,omitempty"`  // name space is like this action.<plugin>.<objectUUID>.<commandName>
	ObjectUUID string `json:"objectUUID,omitempty"` // or use UUID
	Key        string `json:"key,omitempty"`
	Body       any    `json:"body,omitempty"`
}

type Action struct {
	Plugin      string
	ObjectName  string
	CommandName string
	Body        any
}

func CommandParse(input string) *Action {
	parts := strings.Split(input, ".")
	// Ensure we have enough parts to construct the struct
	if len(parts) != 4 {
		return nil
	}
	return &Action{
		Plugin:      parts[0],
		ObjectName:  parts[1],
		CommandName: parts[2],
		Body:        parts[3],
	}
}

type DataPayload struct {
	Data *priority.Value `json:"data"`
	//Connections []*Connection  `json:"connections,omitempty"`
	Topic      string `json:"topic"`
	PortID     string `json:"portID"`
	ObjectUUID string `json:"objectUUID,omitempty"`
	ObjectID   string `json:"objectID,omitempty"`
}

type EventBusPayload struct {
	ObjectUUID                 string        `json:"objectUUID,omitempty"`
	ObjectID                   string        `json:"objectID,omitempty"`
	Topic                      string        `json:"topic,omitempty"`
	ResponseTopic              string        `json:"responseTopic,omitempty"`
	UnsubscribeOnResponseTopic bool          `json:"unsubscribe,omitempty"` // used for when we want to use the EventBus PublishWait and we unsubscribe to the ResponseTopic
	Timeout                    time.Duration `json:"timeout,omitempty"`
	ExpectedData               string        `json:"expectedData,omitempty"` // make it easy for an Obj to decode in incoming data; eg string, map[], user
	Payload                    any           `json:"payload,omitempty"`
}

type Mapping struct {
	ManagerUUID  string       `json:"managerUUID,omitempty"`
	NetworkUUID  string       `json:"networkUUID,omitempty"`
	MapperUUID   string       `json:"mapperUUID,omitempty"`
	Data         any          `json:"data,omitempty"`
	DataPayload  *DataPayload `json:"dataPayload"`
	ExpectedData string       `json:"expectedData,omitempty"` // make it easy for an Obj to decode in incoming data; eg string, map[], user
}

type PrimitivesPayload struct {
	DataType priority.Type      `json:"dataType,omitempty"`
	Priority *priority.Priority `json:"priority,omitempty"`
	Symbol   *string            `json:"symbol,omitempty"`
}

func UnmarshalPayload(resp any) (*Payload, error) {
	payload := &Payload{}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, payload)
	return payload, err
}

func UnmarshalPayloadByte(resp []byte) (*Payload, error) {
	payload := &Payload{}
	err := json.Unmarshal(resp, payload)
	return payload, err
}

func (p *Payload) GetPayload() *Payload {
	return p
}

func (p *Payload) dataNil() {
	if p.DataPayload == nil {
		p.DataPayload = &DataPayload{}
	}

}

func (p *Payload) GetData() *priority.Value {
	return p.GetDataPayload().Data
}

func (p *Payload) SetPriorityData(d *priority.Value) *Payload {
	p.dataNil()
	p.DataPayload.Data = d
	return p
}

func (p *Payload) SetDataPayloadDetails(objectID, objectUUID, portID, topic string) *Payload {
	p.dataNil()
	p.DataPayload.PortID = portID
	p.DataPayload.ObjectID = objectID
	p.DataPayload.ObjectUUID = objectUUID
	p.DataPayload.Topic = topic

	return p
}

func (p *Payload) GetDataPayload() *DataPayload {
	p.dataNil()
	return p.DataPayload
}

func (p *Payload) GetTopic() string {
	p.dataNil()
	return p.DataPayload.Topic
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

func (p *Payload) GetMapping() *Mapping {
	if p.IsMappingNil() {
		return nil
	}
	return p.Mapping
}

// ----------------EVENTBUS------------------

func NewMapping(m *Mapping) *Mapping {
	return m
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
