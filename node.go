package rxlib

import (
	"github.com/NubeIO/schema"
)

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

	BusChannel(inputID string) (chan *Message, bool)
	MessageBus() map[string]chan *Message
	PublishMessage(port *Port, setLastValue ...bool)

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
	BuildSchema(schema *schema.Generated)

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
