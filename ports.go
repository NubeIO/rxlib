package rxlib

import "time"

// type T any
type Port[T any] struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	UUID               string           `json:"uuid"`
	Value              T                `json:"value,omitempty"` // the value after it's had some transformations
	LastUpdated        *time.Time       `json:"lastUpdated"`
	ValueRaw           T                `json:"valueRaw,omitempty"`           // the value before any transformations
	ValueDisplay       any              `json:"valueDisplay"`                 // for example 22 %
	LastUpdatedDisplay string           `json:"lastUpdatedDisplay,omitempty"` // last time it got a message
	Direction          PortDirection    `json:"direction"`                    // input or output
	DataType           PortDataType     `json:"dataType"`                     // float, bool, string, any, json
	Transformations    *Transformations `json:"transformations,omitempty"`    // (if a transformations are used we would add a few extra outputs for valueDisplay and valueRaw)
	// these are optional and used if you want to keep the last value for later use
	PreviousValueSet bool           `json:"-"`
	PreviousValue    *PreviousValue `json:"previousValue,omitempty"`
	// is a value written from another object
	WrittenValueSet          bool          `json:"-"`
	WrittenValue             *WrittenValue `json:"writtenValue,omitempty"`
	AllowMultipleConnections bool          `json:"allowMultipleConnections,omitempty"`

	// port position is where to show the order on the object and where to hide the port or not
	DefaultPosition   int  `json:"defaultPosition"`
	Hide              bool `json:"hide,omitempty"`
	HiddenByDefault   bool `json:"hiddenByDefault,omitempty"`
	OverPositionValue int  `json:"overPositionValue,omitempty"`

	OnMessage func(msg *Payload) // used for the evntbus
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

func (p *Port[T]) SetValue(value T) {
	p.Value = value
}
func (p *Port[T]) SetValueRaw(value T) {
	p.ValueRaw = value
}

func CreatePort[T any](id, name, uuid string, dataType PortDataType, direction PortDirection, f func(message *Payload)) *Port[T] {
	p := &Port[T]{
		ID:        id,
		Name:      name,
		UUID:      uuid,
		DataType:  dataType,
		Direction: direction,
		OnMessage: f,
	}
	return p
}

//func CreatePort[T any](port *Port[T]) *Port[T] {
//	return &Port[T]{
//		ID:       port.ID,
//		Name:     port.Name,
//		UUID:     port.UUID,
//		DataType: port.DataType,
//	}
//}

// ToTime returns a pointer to the passed time.Time value.
func toTime(t time.Time) *time.Time {
	return &t
}

func NewPortFloatCallBack(id string, f func(message *Payload)) *NewPort {
	p := &NewPort{
		ID:        id,
		Name:      id,
		DataType:  PortTypeFloat,
		OnMessage: f,
	}
	//p.DefaultPosition = portOpts(opts...).DefaultPosition
	//p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	//p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortFloat(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: PortTypeFloat,
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
		DataType: PortTypeBool,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortAny(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: PortTypeAny,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

type PreviousValue struct {
	PreviousValue          any       `json:"previousValue,omitempty"`
	PreviousValueRaw       any       `json:"previousValueRaw,omitempty"`
	PreviousValueTimestamp time.Time `json:"previousValueTimestamp,omitempty"`
}

type WrittenValue struct {
	FromUUID   string    `json:"fromUUID"`
	FromPortID string    `json:"fromPortID"`
	Value      any       `json:"previousValue,omitempty"`
	ValueRaw   any       `json:"previousValueRaw,omitempty"`
	Timestamp  time.Time `json:"previousValueTimestamp,omitempty"`
}

type NewPort struct {
	ID                       string
	Name                     string
	DataType                 PortDataType
	AllowMultipleConnections bool `json:"allowMultipleConnections,omitempty"`
	DefaultPosition          int  `json:"defaultPosition"`
	HiddenByDefault          bool `json:"hiddenByDefault,omitempty"`
	OnMessage                func(msg *Payload)
}

type Override struct {
	Value any `json:"value"`
}

type PortFormatString struct {
	ErrorOnMixMax    bool
	MinLengthString  *int     `json:"minLengthString"`
	MaxLengthString  *int     `json:"maxLengthString"`
	AllowEmptyString bool     `json:"allowEmptyString,omitempty"`
	RestrictString   *float64 `json:"restrictString"` // for example don't allow # on an mqtt topic
}

type PortDataType string

const (
	PortTypeJSON   PortDataType = "json"
	PortTypeAny    PortDataType = "any"
	PortTypeFloat  PortDataType = "float"
	PortTypeString PortDataType = "string"
	PortTypeBool   PortDataType = "bool"
)

// some commonly used output names
const (
	OutputName      string = "output"
	OutputErrorName string = "error"
)

// some commonly used input names
const (
	InputName string = "input"
	In1Name   string = "in-1"
)

type FlowDirection string

const (
	DirectionSubscriber FlowDirection = "subscriber"
	DirectionPublisher  FlowDirection = "publisher"
)

type PortDirection string

const (
	Input  PortDirection = "input"
	Output PortDirection = "output"
)
