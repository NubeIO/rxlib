package rxlib

import (
	"github.com/NubeIO/rxlib/libs/bus"
	"github.com/NubeIO/rxlib/payload"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/NubeIO/schema"
	"github.com/patrickmn/go-cache"
	"time"
)

type Object interface {
	New(object Object, opts ...any) Object

	Init() error
	Start() error
	// used normally for the Start() to set it that it has booted
	SetLoaded()
	IsNotLoaded() bool
	IsLoaded() bool // where the Obj Start() method has been called
	CommandList() []*Invoke
	Process() error
	Reset() error // for example this can be called on the 2nd deployment of a add Obj, and we want to reset the count back to zero
	AllowsReset() bool
	Delete() error
	Lock()
	Unlock()
	IsLocked() bool
	IsUnlocked() bool
	InvokePayload(p *payload.Payload) error
	Invoke(command *ExtendedCommand) (*CommandResponse, error)
	CommandObject(command *ExtendedCommand) *CommandResponse   // normally used for objectA to invoke objectB (a way for objects to talk rather than using the eventbus)
	Command(command *ExtendedCommand) *runtime.CommandResponse // normally used for objectA to invoke objectB (a way for objects to talk rather than using the eventbus)
	CommandResponse(response *runtime.CommandResponse)
	AddRuntime(r Runtime)
	Runtime() Runtime
	RemoveObjectFromRuntime()

	GetParentObject() Object
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

	// ports
	NewPort(port *Port)

	NewInputPort(port *NewPort) error
	NewInputPorts(port []*NewPort) error
	NewOutputPort(port *NewPort) error
	NewOutputPorts(port []*NewPort) error
	GetAllPorts() []*Port
	GetPortPayload(portID string) (*payload.Payload, error)
	// SetPortPayload(portID string, data *payload.Payload) error
	GetPortValue(portID string) *runtime.PortValue
	EnablePort(portID string) error
	DisablePort(portID string) error
	IsPortDisable(portID string) (bool, error)
	AddAllTransformations(inputs, outputs []*Port) []error
	OverrideValue(value any, portID string) error
	ReleaseOverride(portID string) error

	CreateConnection(connection *runtime.Connection) // CreateConnection is for just adding a rubix without adding it to the eventbus
	NewOutputConnection(portID, targetUUID, targetPort string) error

	GetPortConnections(portID string) []*runtime.Connection
	GetConnection(uuid string) *runtime.Connection
	GetExistingConnection(sourceObjectUUID, targetObjectUUID, targetPortID string) *runtime.Connection
	GetConnections() []*runtime.Connection
	PortHasConnection(portID string) (has bool, count int)
	RemoveConnection(connection *runtime.Connection) error
	DropConnections() []error
	RemoveOldConnections(newConnections []*runtime.Connection) []error
	AddSubscriptionConnection(sourceObjectUUID, sourcePortID, targetObjectUUID, targetPortID string)

	// inputs
	GetInput(id string) *Port
	InputExists(id string) error
	GetInputs() []*Port
	GetInputByConnection(sourceObjectUUID, outputPortID string) *Port
	GetInputByConnections(sourceObjectUUID, outputPortID string) []*Port
	UpdateInputsValue(portID string, payload *payload.Payload) []error
	SetDynamicInputsCount(count int)
	GetDynamicInputsCount() int
	DeleteInput(id string) error
	BuildDynamicInputs(callback func(portID string, message *payload.Payload)) error
	//SetHasMultipleInputConnections() // set if an output from an object has multiple connection to this objects input or inputs
	//GetHasMultipleInputConnections() map[string]MultipleConnection

	GetOutputs() []*Port
	GetOutput(id string) *Port
	OutputExists(id string) error
	SetOutput(portID string, value any) error // Set current output value & send over the eventbus
	DeleteOutput(id string) error
	// PublishValue eventbus
	PublishValue(portID string) error              // send current port value over the eventbus
	PublishCommand(command *ExtendedCommand) error // send a command over the eventbus, this is used to send external commands (eg; over mqtt)
	Subscribe(topic, handlerID string, callBack func(topic string, e bus.Event))
	SubscribePayload(topic, handlerID string, opts *EventbusOpts, callBack func(topic string, p *payload.Payload, err error))

	// GetRootObject Obj tree
	GetRootObject(uuid string) (Object, error)
	PrintObjectTree(objects map[string]Object)
	GetCompleteChain(objects map[string]Object, uuid string) Chain

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
	GetStats() *runtime.ObjectStats

	// SetInfo -------------------OBJECT INFO------------------
	SetInfo(info *runtime.Info)
	GetInfo() *runtime.Info

	// id
	GetID() string

	// GetObjectType type is for example a driver, service, logic
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
	GetWorkingGroupParent() string
	GetWorkingGroupLeader() string
	GetWorkingGroupLeaderObjectUUID() string

	// plugin
	GetPluginName() string

	// GetMustLiveInObjectType these are needed to know where a know will site in the sidebar in the UI
	GetMustLiveInObjectType() bool
	GetMustLiveParent() bool
	GetRequiresLogger() bool
	AddLogger(trace *Logger)
	Logger() (*Logger, error)
	GetLoggerInfo() ([]string, error)

	// scheam
	GetSchema() *schema.Generated
	// settings
	GetSettings() *runtime.Settings
	SetSettings(settings *runtime.Settings) error

	// GetMeta  meta will also set the Obj-name at parentUUID
	GetMeta() *runtime.Meta
	SetMeta(meta *runtime.Meta) error

	// permissions
	GetPermissions() *runtime.Permissions

	// requirements
	GetRequirements() *runtime.Requirements

	// tags
	AddTag(tag string)
	AddTags(tags ...string)
	GetTag(key string) string
	HasTag(key string) bool
	GetTags() []string
	AddMetaTags(key, value string)
	GetMetaTag(key string) string
	GetMetaTags() map[string]string
	HasMetaTag(key string) bool
	HasMetaTagValue(key, value string) bool

	AddFlag(key, value string)
	GetFlag(key string) string
	GetFlags() map[string]string
	HasFlag(key string) bool
	HasFlagValue(key, value string) bool

	SetCache(key string, data any, expiration time.Duration, overwriteExisting bool) error
	GetCache(key string) (data any, found bool)
	CacheAll() map[string]cache.Item

	GetHelp() string
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

func ConvertObjects(objects []Object) []*runtime.ObjectConfig {
	var out []*runtime.ObjectConfig
	for _, object := range objects {
		out = append(out, ConvertObject(object))
	}
	return out
}

func ConvertObject(object Object) *runtime.ObjectConfig {
	return &runtime.ObjectConfig{
		Id:          object.GetInfo().ObjectID,
		Info:        object.GetInfo(),
		Inputs:      PortsToProto(object.GetInputs()),
		Outputs:     PortsToProto(object.GetOutputs()),
		Meta:        nil,
		Stats:       nil,
		Connections: nil,
		Settings:    nil,
	}
}

func NewMeta(meta *runtime.Meta) *runtime.Meta {
	return meta
}

type Chain struct {
	RootTreeUUIDs       []string
	RootTreeNames       []string
	DescendantTreeUUIDs []string
	DescendantTreeNames []string
}
