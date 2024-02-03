package rxlib

import "time"

type Payload struct {
	Port       *Port       `json:"port"`
	Connection *Connection `json:"connection,omitempty"`
	ObjectUUID string      `json:"objectUUID"`
	ObjectID   string      `json:"objectID"`

	// used for mapping
	MappingUUID       string `json:"mappingUUID,omitempty"`
	RemoteMappingUUID string `json:"remoteMappingUUID,omitempty"`
}

func (p *Payload) GetValue() any {
	return p.Port.Value
}

type PayloadValue struct {
	Value     any `json:"value"`
	Timestamp *time.Time
}

func (p *Payload) GetValueLastUpdated() *PayloadValue {
	return &PayloadValue{
		Value:     p.GetValue(),
		Timestamp: p.GetLastUpdated(),
	}
}

func (p *Payload) GetLastUpdated() *time.Time {
	return p.Port.LastUpdated
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

func (p *Payload) SetMappingUUID(value string) {
	p.MappingUUID = value
}

func (p *Payload) GetMappingUUID() string {
	return p.MappingUUID
}

func (p *Payload) GetRemoteMappingUUID() string {
	return p.RemoteMappingUUID
}

func (p *Payload) SetRemoteMappingUUID(value string) {
	p.RemoteMappingUUID = value
}
