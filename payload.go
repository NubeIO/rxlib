package rxlib

import "time"

type Payload struct {
	Port       *Port[any] `json:"port"`
	ObjectUUID string     `json:"objectUUID"`
	ObjectID   string     `json:"objectID"`
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

func (p *Payload) GetPort() *Port[any] {
	return p.Port
}

func (p *Payload) FromObject() string {
	return p.ObjectUUID
}

func (p *Payload) FromObjectID() string {
	return p.ObjectID
}
