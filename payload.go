package rxlib

import (
	"github.com/NubeIO/rxlib/priority"
	"time"
)

func NewPayload(body *Payload) *Payload {
	return body
}

type Payload struct {
	Port       *Port       `json:"port"`
	Connection *Connection `json:"connection,omitempty"`
	ObjectUUID string      `json:"objectUUID"`
	ObjectID   string      `json:"objectID"`
	// used for mapping
	Mapping *Mapping `json:"mapping,omitempty"`

	ExpectedData string `json:"expectedData"` // make it easy for an object to decode in incoming data; eg string, map[], user
	Data         any    `json:"data"`
}

type Mapping struct {
	ManagerUUID       string            `json:"managerUUID,omitempty"`
	NetworkUUID       string            `json:"networkUUID,omitempty"`
	MapperUUID        string            `json:"mapperUUID,omitempty"`
	Data              any               `json:"data,omitempty"`
	PrimitivesPayload PrimitivesPayload `json:"primitivesPayload,omitempty"`
}

type PrimitivesPayload struct {
	DataType priority.Type      `json:"dataType,omitempty"`
	Priority *priority.Priority `json:"priority,omitempty"`
	Symbol   *string            `json:"symbol,omitempty"`
}

func (p *Payload) GetValue() any {
	return p.Port.Data
}

func (p *Payload) GetDataPriority() *priority.Priority {
	if p.Port.DataPriority == nil {
		return nil
	}
	if p.Port.DataPriority.Priority == nil {
		return nil
	}
	return p.Port.DataPriority.Priority
}

type PayloadValue struct {
	Value     any `json:"value"`
	Timestamp *time.Time
}

func (p *Payload) GetPortName() string {
	return p.Port.Name
}

func (p *Payload) GetPort() *Port {
	return p.Port
}

func (p *Payload) FromObject() string {
	return p.ObjectUUID
}

func (p *Payload) FromObjectID() string {
	return p.ObjectID
}

func (p *Payload) GetConnection() *Connection {
	return p.Connection
}

func (p *Payload) SetObjectUUID(value string) {
	p.ObjectUUID = value
}
