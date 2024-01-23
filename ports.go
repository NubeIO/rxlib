package rxlib

type Port struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Value           any              `json:"value,omitempty"`       // the value after it's had some transformations
	ValueRaw        any              `json:"valueRaw,omitempty"`    // the value before any transformations
	ValueDisplay    any              `json:"valueDisplay"`          // for example 22 %
	LastUpdated     string           `json:"lastUpdated,omitempty"` // last time it got a message
	Direction       PortDirection    `json:"direction"`
	DataType        PortDataType     `json:"dataType"`
	Position        int              `json:"position"`                  // node position to display in the UI
	Transformations *Transformations `json:"transformations,omitempty"` // (if a transformations are used we would add a few extra outputs for valueDisplay and valueRaw)
}

type NewPort struct {
	ID       string
	Name     string
	DataType PortDataType
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
