package rxlib

import (
	"github.com/NubeIO/schema"
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
	SetLoaded()
	IsLoaded() bool // where the object Start() method has been called
	Process() error
	Reset() error // for example this can be called on the 2nd deploy of a counter object, and we want to reset the count back to zero
	AllowsReset() bool
	Delete() error
	Lock()
	Unlock()
	IsLocked() bool
	IsUnlocked() bool

	// runtime objects
	AddRuntimeToObject(runtimeObjects map[string]Object) // gives each object access to every other object
	GetRuntimeObjects() map[string]Object
	GetForeignObject(objectUUID string) (obj Object, exists bool)
	CheckForeignObjectOutputExists(objectUUID, portID string) (*Port, error)
	CheckForeignObjectInputExists(objectUUID, portID string) (*Port, error)
	GetObjectsByType(objectID string) []Object // for example get all math/add Object
	RemoveObjectFromRuntime()
	GetChildObjects() []Object                      // get all the object inside a folder
	GetChildObjectsByType(objectID string) []Object // for example get all modbus/device that are inside its parent modbus/network Object
	GetParentObject(uuid string) (obj Object, exists bool)
	GetParentUUID() string

	// AddExtension extension are a way to extend the functionalists of an object; for example add a history extension
	AddExtension(extension Object)
	GetExtensions() map[string]Object
	GetExtension(id string) Object
	DeleteExtension(name string)

	// ports
	NewPort(port *Port)
	NewInputPort(port *NewPort) error
	NewInputPorts(port []*NewPort) error
	NewOutputPort(port *NewPort) error
	NewOutputPorts(port []*NewPort) error
	GetAllPorts() []*Port
	// connections
	AddConnection(connection *Connection) error
	GetConnection(uuid string) (*Connection, error)
	GetConnections() []*Connection
	UpdateConnections(connections []*Connection) *UpdateConnectionsReport
	RemoveConnection(connection *Connection) *RemoveConnectionReport
	RemoveAllConnections() []*RemoveConnectionReport

	// inputs
	GetInput(id string) *Port
	GetInputs() []*Port
	SetInputValue(id string, value any) error

	// ouputs
	GetOutputs() []*Port
	GetOutput(id string) *Port
	// WriteValue update the port value; pass in option withTimestamp to timestamp to write
	WriteValue(portID string, value any, withTimestamp ...bool) error

	//GetAllObjectValues ObjectValue are a way for one node to direly get and send data to another node
	// PreviousValue is the last value saved
	// WrittenValue is a value written from another object; this is useful for example on network object where the network is doing the polling and can quickly update the devices/points
	GetAllObjectValues() []*ObjectValue
	SetOutputPreviousValue(id string, value *PreviousValue) error
	GetOutputPreviousValue(id string) *PreviousValue
	SetInputPreviousValue(id string, value *PreviousValue) error
	GetInputPreviousValue(id string) *PreviousValue

	SetOutputWrittenValue(id string, value *WrittenValue) error
	GetOutputWrittenValue(id string) *WrittenValue
	SetInputWrittenValue(id string, value *WrittenValue) error
	GetInputWrittenValue(id string) *WrittenValue

	// data TODO maybe add a cache timeout, also a GetTheDelete() and a Delete()
	GetData() map[string]any
	SetData(key string, data any) // addData is a way for a node to store something in memory
	DeleteData(key string) error
	GetDataByKey(key string, out any) error

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
	SetLoopCount(count uint)
	GetLoopCount() uint
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
	AddLogger(trace *Logger)
	Logger() (*Logger, error)
	GetLoggerInfo() (*LoggerOpts, error)

	// scheam
	GetSchema() *schema.Generated
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
