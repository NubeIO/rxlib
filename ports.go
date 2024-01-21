package rxlib

type Port struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Value          any           `json:"value,omitempty"`
	ValueFormatted any           `json:"valueFormatted"`        // for example 22 %
	LastUpdated    string        `json:"lastUpdated,omitempty"` // last time it got a message
	Direction      PortDirection `json:"direction"`
	DataType       PortDataType  `json:"dataType"`
	FormatNumber   *FormatNumber `json:"portDataAttributes,omitempty"`
	PermitNull     bool          `json:"permitNull"` // if true will set the default value of golang types;
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
