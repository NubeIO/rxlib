package rxlib

import (
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/priority"
)

type Port struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	UUID string `json:"uuid"`

	// Input/Output port values
	Data         *priority.DataValue    `json:"value,omitempty"`     // value should be used for anything
	DataPriority *priority.DataPriority `json:"dataFloat,omitempty"` // values for float
	primitives   *priority.Primitives

	Direction PortDirection `json:"direction"` // input or output
	DataType  priority.Type `json:"dataType"`  // float, bool, string, any, json

	// these are optional and used if you want to keep the last value for later use
	PreviousValueSet bool                    `json:"-"`
	PreviousValue    *priority.PreviousValue `json:"previousValue,omitempty"`
	// is a value written from another Obj
	WrittenValueSet          bool                   `json:"-"`
	WrittenValue             *priority.WrittenValue `json:"writtenValue,omitempty"`
	AllowMultipleConnections bool                   `json:"allowMultipleConnections,omitempty"`

	// port position is where to show the order on the Obj and where to hide the port or not
	DefaultPosition   int  `json:"defaultPosition"`
	Hide              bool `json:"hide,omitempty"`
	HiddenByDefault   bool `json:"hiddenByDefault,omitempty"`
	OverPositionValue int  `json:"overPositionValue,omitempty"`

	OnMessage func(msg *Payload) `json:"-"` // used for the evntbus

}

func (p *Port) GetID() string {
	return p.ID
}

func (p *Port) GetUUID() string {
	return p.UUID
}

func (p *Port) GetName() string {
	return p.Name
}

func (p *Port) SetName(v string) string {
	p.Name = v
	return p.Name
}

// SetData to be used for anything but an int, float, bool
func (p *Port) SetData(value any) {
	if p.Data == nil {
		p.Data = &priority.DataValue{}
	}
	p.Data.Value = value
}

func (p *Port) GetData() any {
	return p.Data
}

//-----------------------Priority-----------------------

func (p *Port) initDataPriority(body *priority.NewPrimitiveValue) error {
	pri, prim, err := priority.NewPrimitive(body)
	p.primitives = prim
	if err != nil {
		return err
	}
	p.DataPriority = pri
	return nil
}

func (p *Port) InitDataPriorityFloat(body *priority.NewPrimitiveValue) error {
	return p.initDataPriority(body)
}

func (p *Port) SetPriorityNull(priorityNumber int) {
	p.DataPriority.Priority.SetNull(priorityNumber)
}

func (p *Port) WriteFloat(value float64) error {
	if p.DataPriority == nil {
		panic("rxlib.SetDataPriority DataPriority can not be empty, please InitDataPriority() first")
	}
	result, err := p.primitives.UpdateValueFloat(value)
	if err != nil {
		return err
	}
	p.DataPriority = result
	return nil
}

func (p *Port) OverrideWriteFloatPriority(value float64, priority int) error {
	if p.DataPriority == nil {
		panic("rxlib.SetDataPriority DataPriority can not be empty, please InitDataPriority() first")
	}
	result, err := p.primitives.UpdateValueAndGenerateResult(nil, 0, nils.ToFloat64(value), priority)
	if err != nil {
		return err
	}
	p.DataPriority = result
	return nil
}

func (p *Port) SetWrittenValue(v *priority.WrittenValue) {
	p.WrittenValueSet = true
	p.WrittenValue = v

}

func (p *Port) GetWrittenValue() *priority.WrittenValue {
	return p.WrittenValue
}

func (p *Port) GetWrittenValueCurrent() (value any, ok bool) {
	if p.WrittenValue != nil {
		return p.WrittenValue.Value, true
	}
	return nil, false
}

type PortOpts struct {
	DefaultPosition          int  `json:"defaultPosition"`
	HiddenByDefault          bool `json:"hiddenByDefault,omitempty"`
	AllowMultipleConnections bool `json:"allowMultipleConnections,omitempty"`
}

func portOpts(opts ...*PortOpts) *PortOpts {
	p := &PortOpts{}
	if len(opts) > 0 {
		if opts[0] != nil {
			p.DefaultPosition = opts[0].DefaultPosition
			p.HiddenByDefault = opts[0].HiddenByDefault
			p.AllowMultipleConnections = opts[0].AllowMultipleConnections
		}
	}
	return p
}

func NewPortFloatCallBack(id string, f func(message *Payload), opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:        id,
		Name:      id,
		DataType:  priority.TypeFloat,
		OnMessage: f,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortFloat(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeFloat,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortBool(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeBool,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

type NewPort struct {
	ID                       string
	Name                     string
	DataType                 priority.Type
	AllowMultipleConnections bool               `json:"allowMultipleConnections,omitempty"`
	DefaultPosition          int                `json:"defaultPosition"`
	HiddenByDefault          bool               `json:"hiddenByDefault,omitempty"`
	OnMessage                func(msg *Payload) `json:"-"`
}

type Override struct {
	Value any `json:"value"`
}

type PortFormatString struct {
	ErrorOnMixMax    bool     `json:"errorOnMixMax"`
	MinLengthString  *int     `json:"minLengthString"`
	MaxLengthString  *int     `json:"maxLengthString"`
	AllowEmptyString bool     `json:"allowEmptyString,omitempty"`
	RestrictString   *float64 `json:"restrictString"` // for example don't allow # on an mqtt topic
}

// some commonly used output names
const (
	OutputName      string = "output"
	OutputErrorName string = "Err"
)

// some commonly used input names
const (
	InputName string = "input"
	In1Name   string = "in-1"
	In2Name   string = "in-2"
)

type FlowDirection string

const (
	DirectionSubscriber      FlowDirection = "subscriber"
	DirectionPublisher       FlowDirection = "publisher"
	DirectionRequestResponse FlowDirection = "request-response"
)

type PortDirection string

const (
	Input  PortDirection = "input"
	Output PortDirection = "output"
)
