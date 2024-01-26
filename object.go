package rxlib

import (
	"github.com/NubeIO/schema"
	"github.com/NubeIO/tracer"
	"github.com/gin-gonic/gin"
)

type Chain struct {
	RootTreeUUIDs       []string
	RootTreeNames       []string
	DescendantTreeUUIDs []string
	DescendantTreeNames []string
}

type Object interface {
	New(object Object, opts ...any) Object

	// Start the processing
	Start() error
	Process() error
	Delete() error
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	// runtime objects
	AddRuntimeToObject(runtimeObjects map[string]Object) // gives each object access to every other object
	GetRuntimeObjects() map[string]Object
	GetObject(uuid string) (obj Object, exists bool)
	GetObjectsByType(objectID string) []Object // for example get all math/add Object
	RemoveObjectFromRuntime()
	GetChildObjects() []Object                      // get all the object inside a folder
	GetChildObjectsByType(objectID string) []Object // for example get all modbus/device that are inside its parent modbus/network Object
	GetParentObject(uuid string) (obj Object, exists bool)
	GetParentUUID() string

	// event bus
	BusChannel(inputID string) (chan *Message, bool)
	MessageBus() map[string]chan *Message
	PublishMessage(port *Port)

	// ports
	NewPort(port *Port)
	NewInputPort(port *NewPort) error
	NewInputPorts(port []*NewPort) error
	NewOutputPort(port *NewPort) error
	NewOutputPorts(port []*NewPort) error
	GetAllPorts() []*Port
	// connections
	AddConnection(connection *Connection)
	GetConnections() []*Connection
	UpdateConnections(connections []*Connection)

	// inputs
	GetInput(id string) *Port
	GetInputs() []*Port
	SetInputValue(id string, value any) error

	// ouputs
	GetOutputs() []*Port
	GetOutput(id string) *Port

	// values
	GetAllObjectValues() []*ObjectValue
	SetOutputPreviousValue(id string, value *PreviousValue) error
	GetOutputPreviousValue(id string) *PreviousValue
	SetInputPreviousValue(id string, value *PreviousValue) error
	GetInputPreviousValue(id string) *PreviousValue

	SetOutputWrittenValue(id string, value *WrittenValue) error
	GetOutputWrittenValue(id string) *WrittenValue
	SetInputWrittenValue(id string, value *WrittenValue) error
	GetInputWrittenValue(id string) *WrittenValue

	// scheam
	GetSchema() *schema.Generated

	// data TODO maybe add a cache timeout, also a GetTheDelete() and a Delete()
	GetData() map[string]any
	AddData(key string, data any) // addData is a way for a node to store something in memory
	GetDataByKey(key string, out interface{}) error

	// GetRootObject object tree
	GetRootObject(uuid string) (Object, error)
	PrintObjectTree(objects map[string]Object)
	GetCompleteChain(objects map[string]Object, uuid string) Chain

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

	// category
	GetCategory() string

	// working group; a group of objects that work together like a network driver
	GetWorkingGroup() string
	GetWorkingGroupObjects() []string
	GetWorkingGroupChildObjects() []string
	GetWorkingGroupParent() string
	GetWorkingGroupLeader() string

	// plugin
	GetPluginName() string

	//GetMustLiveInObjectType these are needed to know where a know will site in the sidebar in the UI
	GetMustLiveInObjectType() bool
	GetMustLiveParent() bool
	GetRequiresLogger() bool
	AddLogger(trace *tracer.Logger)
	Logger() (*tracer.Logger, error)

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
