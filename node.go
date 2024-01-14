package rxlib

import (
	"github.com/NubeIO/schema"
)

type portDataType string

const (
	portTypeAny    portDataType = "any"
	portTypeFloat  portDataType = "float"
	portTypeString portDataType = "string"
	portTypeBool   portDataType = "bool"
)

type flowDirection string

const (
	DirectionSubscriber flowDirection = "subscriber"
	DirectionPublisher  flowDirection = "publisher"
)

type portDirection string

const (
	input  portDirection = "input"
	output portDirection = "output"
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
	Direction   portDirection `json:"direction"`
	DataType    portDataType  `json:"dataType"`
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
	FlowDirection flowDirection `json:"flowDirection"` // subscriber is if it's in an input and publisher if It's for an output

}

type Node interface {
	New(nodeUUID, name string, bus *EventBus, settings *Settings) Node
	Start()
	Delete()
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	// node details/info
	GetUUID() string
	GetParentUUID() string
	GetPluginName() string
	GetApplicationUse() string
	GetID() string
	GetNodeName() string

	// ports
	NewPort(port *Port)

	// connections
	AddConnection(connection *Connection)
	GetConnections() []*Connection
	UpdateConnections(connections []*Connection)

	// inputs
	GetInput(id string) *Port
	GetInputs() []*Port
	SetInputValue(id string, value interface{})

	// ouputs
	GetOutputs() []*Port

	// values
	GetAllNodeValues() []*NodeValue
	GetAllPortValues() []*Port
	GetAllInputValues() []*Port
	GetAllOutputValues() []*Port
	SetLastValue(port *Port)
	GetPortValue(portID string) (*Port, error)

	// scheam
	GetSchema() *schema.Generated
	AddSchema()

	// details
	SetDetails(details *Details)
	GetDetails() *Details

	// data
	AddData(key string, data any)
	GetDataByKey(key string, out interface{}) error

	// meta
	GetData() map[string]any
	setMeta(opts *Options)
	GetMeta() *Meta

	// settings
	AddSettings(settings *Settings)
	GetSettings() *Settings
	UpdateSettings(settings *Settings)

	// runtime nodes
	AddRuntime(runtimeNodes map[string]Node) // gives each node access to every other node
	GetRuntimeNodes() map[string]Node
	AddToNodeToRuntime(node Node) Node
	RemoveNodeFromRuntime()

	// child nodes
	RegisterChildNode(child Node)
	GetChildNodes() []Node
	GetChildNode(uuid string) Node
	GetChildsByType(nodeID string) []Node
	GetPortValuesChildNode(uuid string) []*Port
	SetLastValueChildNode(uuid string, port *Port)

	// options
	AddOptions(opts *Options)
	GetOptions() *Options

	// wants Services
	AddServices()
	SupportsServices() bool
}

type NodeValue struct {
	NodeId   string  `json:"nodeId"`
	NodeUUID string  `json:"nodeUUID"`
	Ports    []*Port `json:"ports"`
}

func NewOptions(opts *Options) *Options {
	return opts
}

type Options struct {
	addToNodesMap bool
	Meta          *Meta `json:"meta"`
}

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	Position   Position `json:"position"`
	ParentUUID string   `json:"parentUUID"`
}
