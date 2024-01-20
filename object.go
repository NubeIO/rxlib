package rxlib

import (
	"github.com/NubeIO/schema"
	"github.com/gin-gonic/gin"
)

type Object interface {
	New(object Object, opts ...any) Object

	// Start the processing
	Start() error
	Delete() error
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	BusChannel(inputID string) (chan *Message, bool)
	MessageBus() map[string]chan *Message
	PublishMessage(port *Port, setLastValue ...bool)

	// ports
	NewPort(port *Port)
	NewInputPort(id, name string, dataType PortDataType)
	NewOutputPort(id, name string, dataType PortDataType)

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
	GetOutput(id string) *Port

	// values
	GetAllObjectValues() []*ObjectValue
	GetAllPortValues() []*Port
	GetAllInputValues() []*Port
	GetAllOutputValues() []*Port
	SetLastValue(port *Port)
	GetPortValue(portID string) (*Port, error)

	// scheam
	GetSchema() *schema.Generated

	// data TODO maybe add a cache timeout, also a GetTheDelete() and a Delete()
	GetData() map[string]any
	AddData(key string, data any) // addData is a way for a node to store something in memory
	GetDataByKey(key string, out interface{}) error

	// runtime objects
	AddRuntimeToObject(runtimeObjects map[string]Object) // gives each object access to every other object
	GetRuntimeObjects() map[string]Object
	RemoveObjectFromRuntime()
	//AddObjectToRuntime(object Object)

	// child objects
	AddDefinedChildObjects(objectID ...string) // to show the UI a objects childs that are defined by the plugin developer
	GetDefinedChildObjects() []string
	RegisterChildObject(child Object)
	GetChildObjects() []Object
	GetChildObject(uuid string) Object
	DeleteChildObject(uuid string) error
	GetChildsByType(objectID string) []Object
	GetPortValuesChildObject(uuid string) []*Port
	SetLastValueChildObject(uuid string, port *Port)

	// RunValidation -------------------VALIDATION INFO------------------
	// ValidationBuilder validation for example, you want to add a new network so lets run some checks eg; is network interface available
	RunValidation()
	AddValidation(key string)
	DeleteValidation(key string)
	GetValidations() map[string]*ErrorsAndValidation
	GetValidation(key string) (*ErrorsAndValidation, bool)
	SetError(key string, err error)
	SetValidationError(key string, m *ValidationMessage)
	SetHalt(key string, m *ValidationMessage)
	SetValidation(key string, m *ValidationMessage)

	// SetStatus -------------------STATS INFO------------------
	// this is for the node status
	SetStatus(status ObjectStatus)
	SetLoopCount(count int)
	GetLoopCount() int
	IncrementLoopCount()
	ResetLoopCount()
	GetStats() *ObjectStats
	AddCustomStat(name string, stat *CustomStatus)
	GetCustomStat(name string) (*CustomStatus, bool)
	DeleteCustomStat(name string)
	UpdateCustomStat(name string, stat *CustomStatus)

	// SetInfo -------------------OBJECT INFO------------------
	SetInfo(info *Info)
	GetInfo() *Info

	// id
	GetID() string

	// object type is for example a driver, service, logic
	GetObjectType() ObjectType

	// uuid, set from Meta
	GetUUID() string

	// name, set from Meta
	GetName() string

	// parent uuid, set from Meta
	GetParentUUID() string

	// category
	GetCategory() string

	// plugin
	GetPluginName() string

	// settings
	GetSettings() *Settings
	SetSettings(settings *Settings) error

	// GetMeta  meta will also set the object-name at parentUUID
	GetMeta() *Meta
	SetMeta(meta *Meta) error

	// permissions
	GetPermissions() *Permissions

	// requirements
	GetRequirements() *Requirements

	// tags
	AddObjectTags(objectTypeTag ...ObjectTypeTag)
	GetObjectTags() []ObjectTypeTag

	AddRouterGroup(c *gin.RouterGroup)
}

type ObjectValue struct {
	ObjectId   string  `json:"objectID"`
	ObjectUUID string  `json:"objectUUID"`
	Ports      []*Port `json:"ports"`
}

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	ObjectUUID string   `json:"uuid"`                 // comes from UI need to set in objectInfo
	ObjectName string   `json:"name"`                 // comes from UI need to set in objectInfo
	ParentUUID string   `json:"parentUUID,omitempty"` // comes from UI need to set in objectInfo
	Position   Position `json:"position"`
}
