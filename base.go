package rxlib

import (
	"github.com/NubeIO/rxlib/libs/history"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/schema"
	"github.com/gin-gonic/gin"
	"github.com/mustafaturan/bus/v3"
	"time"
)

type Object interface {
	New(object Object, opts ...any) Object

	Init() error
	Start() error
	SetLoaded() // used normally for the Start() to set it that it has booted
	IsNotLoaded() bool
	IsLoaded() bool // where the Obj Start() method has been called
	Invoke(key string, dataType priority.Type, data any) *ObjectCommandResponse
	InvokePayload(key string, dataType priority.Type, payload *Payload) *ObjectCommandResponse
	InvokePayloads(key string, dataType priority.Type, payload []*Payload) *ObjectCommandResponse
	CommandList() []*Invoke
	Process() error
	Reset() error // for example this can be called on the 2nd deployment of a counter Obj, and we want to reset the count back to zero
	AllowsReset() bool
	Delete() error
	Lock()
	Unlock()
	IsLocked() bool
	IsUnlocked() bool

	Command(command *Command) *CommandResponse // normally used for objectA to invoke objectB (a way for objects to talk rather than using the eventbus)

	AddRuntime(r Runtime)
	Runtime() Runtime
	RemoveObjectFromRuntime()

	GetChildObjectsByType(objectID string) []Object // for example get all modbus/device that are inside its parent modbus/network Object
	GetChildObjects() []Object
	GetParentObject() Object
	GetChildObject(objectUUID string) Object
	GetParentUUID() string

	// AddExtension extension are a way to extend the functionalists of an Obj; for example add a history extension
	AddExtension(extension Object) error
	GetExtensions() []Object
	GetExtension(id string) (Object, error)
	DeleteExtension(name string) error

	SetRequiredExtensions(extension []*Extension)
	GetRequiredExtensions() []*Extension
	RequiredExtensionListCount() (extensionsCount int) // get a count if there are any required extensions or not
	IsExtensionsAdded(objectID string) (addedCount int)
	GetRequiredExtensionByName(extensionName string) *Extension
	//HistoryManager history's
	GetHistoryManager() history.Manager
	SetHistoryManager(h history.Manager)

	// ports
	NewPort(port *Port)
	NewInputPort(port *NewPort) error
	NewInputPorts(port []*NewPort) error
	NewOutputPort(port *NewPort) error
	NewOutputPorts(port []*NewPort) error
	GetAllPorts() []*Port
	EnablePort(portID string) error
	DisablePort(portID string) error
	IsPortDisable(portID string) (bool, error)
	AddAllTransformations(inputs, outputs []*Port) []error
	OverrideValue(value any, portID string) error
	ReleaseOverride(portID string) error

	CreateConnection(connection *Connection) // CreateConnection is for just adding a rubix without adding it to the eventbus
	NewOutputConnection(portID, targetUUID, targetPort string) error

	GetConnection(uuid string) *Connection
	GetExistingConnection(sourceObjectUUID, targetObjectUUID, targetPortID string) *Connection
	GetConnections() []*Connection

	RemoveConnection(connection *Connection) error
	DropConnections() []error
	RemoveOldConnections(newConnections []*Connection) []error
	AddSubscriptionConnection(sourceObjectUUID, sourcePortID, targetObjectUUID, targetPortID string)

	// inputs
	GetInput(id string) *Port
	InputExists(id string) error
	GetInputByUUID(uuid string) *Port
	GetInputs() []*Port
	GetInputByConnection(sourceObjectUUID, outputPortID string) *Port
	GetInputByConnections(sourceObjectUUID, outputPortID string) []*Port
	UpdateInputsValues(payload *Payload) []error
	GetInputValue(portID string) *priority.Value
	AllInputsValues() map[string]*priority.Value
	GetInputData(portID string) *priority.PriorityData
	AllInputsData() map[string]*priority.PriorityData

	GetOutputs() []*Port
	GetOutput(id string) *Port
	OutputExists(id string) error
	GetOutputByUUID(uuid string) *Port
	SetOutput(portID string, value any) error // Set current output value & send over the eventbus
	PublishValue(portID string) error         // send current port value over the eventbus
	Subscribe(topic, handlerID string, callBack func(topic string, e bus.Event))
	SubscribePayload(topic, handlerID string, opts *EventbusOpts, callBack func(topic string, p *Payload, err error))
	GetOutputValue(portID string) *priority.Value
	AllOutputsValues() map[string]*priority.Value
	GetOutputData(portID string) *priority.PriorityData
	AllOutputsData() map[string]*priority.PriorityData

	// GetRootObject Obj tree
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

	// Obj type is for example a driver, service, logic
	GetObjectType() ObjectType

	// uuid, set from Meta
	GetUUID() string

	// name, set from Meta
	GetName() string
	SetName(v string) string

	// category
	GetCategory() string

	// working group; a group of objects that work together like a network driver
	GetWorkingGroup() string
	GetWorkingGroupObjects() []string
	GetWorkingGroupChildObjects() []string
	GetWorkingGroupParent() string
	GetWorkingGroupLeader() string
	GetWorkingGroupLeaderObjectUUID() string

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

	// GetMeta  meta will also set the Obj-name at parentUUID
	GetMeta() *Meta
	SetMeta(meta *Meta) error

	// permissions
	GetPermissions() *Permissions

	// requirements
	GetRequirements() *Requirements
	GetRubixServicesRequirement() []*RubixServicesRequirement

	// tags
	AddObjectTags(objectTypeTag ...ObjectTypeTag)
	GetObjectTags() []ObjectTypeTag

	SetCache(key string, data any, expiration time.Duration, overwriteExisting bool) error
	GetCache(key string) (data any, found bool)

	AddRouterGroup(c *gin.RouterGroup)
}

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

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

func NewMeta(meta *Meta) *Meta {
	return meta
}

type Chain struct {
	RootTreeUUIDs       []string
	RootTreeNames       []string
	DescendantTreeUUIDs []string
	DescendantTreeNames []string
}
