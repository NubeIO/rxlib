package rxlib

import (
	"github.com/NubeIO/schema"
)

type Object interface {
	New(objectUUID, name string, bus *EventBus, settings *Settings) Object
	Start()
	Delete()
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	// object details/info
	GetUUID() string
	GetParentUUID() string
	GetPluginName() string
	GetApplicationUse() string
	GetID() string
	GetObjectName() string

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
	GetAllObjectValues() []*ObjectValue
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

	// runtime objects
	AddRuntime(runtimeObjects map[string]Object) // gives each object access to every other object
	GetRuntimeObjects() map[string]Object
	AddToObjectToRuntime(object Object) Object
	RemoveObjectFromRuntime()

	// child objects
	RegisterChildObject(child Object)
	GetChildObjects() []Object
	GetChildObject(uuid string) Object
	GetChildsByType(objectID string) []Object
	GetPortValuesChildObject(uuid string) []*Port
	SetLastValueChildObject(uuid string, port *Port)

	// options
	AddOptions(opts *Options)
	GetOptions() *Options

	// wants Services
	AddServices()
	SupportsServices() bool
}

type ObjectValue struct {
	ObjectId   string  `json:"objectId"`
	ObjectUUID string  `json:"objectUUID"`
	Ports      []*Port `json:"ports"`
}

func NewOptions(opts *Options) *Options {
	return opts
}

type Options struct {
	addToObjectsMap bool
	Meta            *Meta `json:"meta"`
}

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	Position   Position `json:"position"`
	ParentUUID string   `json:"parentUUID"`
}
