package rxlib

type PortDataType string

const (
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

type Details struct {
	Category    string  `json:"category"`
	ParentID    *string `json:"parentID"`
	HasServices bool    `json:"hasServices"`
}

type Port struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Value       interface{}   `json:"value,omitempty"`
	LastUpdated string        `json:"lastUpdated,omitempty"` // last time it got a message
	Direction   PortDirection `json:"direction"`
	DataType    PortDataType  `json:"dataType"`
}

type Settings struct {
	Value interface{} `json:"value"`
}

// Connection defines a structure for input subscriptions.
type Connection struct {
	SourceUUID    string        `json:"source"`
	SourcePort    string        `json:"sourceHandle"`
	TargetUUID    string        `json:"target"`
	TargetPort    string        `json:"targetHandle"`
	FlowDirection FlowDirection `json:"flowDirection"` // subscriber is if it's in an input and publisher if It's for an output

}
