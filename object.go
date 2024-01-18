package rxlib

import (
	"github.com/NubeIO/schema"
	"github.com/gin-gonic/gin"
)

type Info struct {
	ObjectID    string
	ObjectUUID  string
	Name        string
	PluginName  string
	Application string
}

// ObjectInfo
//   - ObjectID
//   - ObjectUUID
//   - Name
//   - PluginName
//   - Application
func ObjectInfo(s ...string) *Info {
	setup := &Info{
		ObjectID:   "",
		ObjectUUID: "",
		Name:       "",
		PluginName: "",
	}

	// Check the length of s and assign values accordingly
	if len(s) > 0 {
		setup.ObjectID = s[0]
	}
	if len(s) > 1 {
		setup.ObjectUUID = s[1]
	}
	if len(s) > 2 {
		setup.Name = s[2]
	}
	if len(s) > 3 {
		setup.PluginName = s[3]
	}
	if len(s) > 4 {
		setup.PluginName = s[4]
	}

	return setup
}

type Object interface {
	New(info *Info, object Object) Object
	New3(objectUUID, name string, bus *EventBus, settings *Settings) Object
	Start()
	Delete()
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	// object details/info
	GetID() string
	SetID(uuid string)
	GetUUID() string
	SetUUID(uuid string)
	GetParentUUID() string
	GetPluginName() string
	GetApplicationUse() string
	GetObjectName() string

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

	// values
	GetAllObjectValues() []*ObjectValue
	GetAllPortValues() []*Port
	GetAllInputValues() []*Port
	GetAllOutputValues() []*Port
	SetLastValue(port *Port)
	GetPortValue(portID string) (*Port, error)

	// scheam
	GetSchema() *schema.Generated

	// settings
	GetSettings() *Settings
	AddUpdateSettings(settings *Settings)

	// details
	SetDetails(details *Details)
	GetDetails() *Details
	AddObjectTypeRequirement(value ObjectTypeRequirements)
	AddObjectTypeRequirements(requirement ...ObjectTypeRequirements) // ObjectTypeRequirements is somthing like an object can only be added once
	GetTypeRequirements() map[ObjectRequirement]ObjectTypeRequirements
	GetTypeRequirement(key ObjectRequirement) ObjectTypeRequirements
	TypeRequirementExists(key ObjectRequirement) bool
	AddObjectTypeTags(objectTypeTag ...ObjectTypeTag)
	GetObjectTypeTags() []ObjectTypeTag

	// data TODO maybe add a cache timeout, also a GetTheDelete() and a Delete()
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

	// options
	AddOptions(opts *Options)
	GetOptions() *Options

	// meta
	GetData() map[string]any
	setMeta(opts *Options)
	GetMeta() *Meta

	// validation
	RunValidation()                                     // for example, you want to add a new network so lets run some checks eg; is network interface available
	GetValidation() map[string]ErrorsAndValidation      // get them
	SetValidation(data map[string]ErrorsAndValidation)  // set them
	AddValidation(key string, data ErrorsAndValidation) // add one
	NewValidation(key string, issue string)
	NewError(key string, err error)
	NewHalt(key, issue, explanation string)
	DeleteValidation(key string) bool
	SetValidationFlag(bool)
	SetErrorFlag(bool)
	SetHaltFlag(bool)     // we may halt/disable the operation of the object execution do a error
	HaltFlag() bool       // for example, we halt the operation for an object as a key requirement has not been filled, for example a database connection could not be made so disable the running of the logic in the object
	ValidationFlag() bool // an error is somthing that is not a validation error, this is somthing that we may not want to show the user
	ErrorFlag() bool
	SetObjectStatus(value ObjectStatus)
	GetObjectStatus() ObjectStatus

	AddRouterGroup(c *gin.RouterGroup)
}

// ErrorsValidation error, validation
type ErrorsValidation string

const (
	TypeHalt       ErrorsValidation = "halt"
	TypeError      ErrorsValidation = "error"
	TypeValidation ErrorsValidation = "validation"
)

type ErrorsAndValidation struct {
	Type              ErrorsValidation `json:"type"`
	Error             error            `json:"-"`
	ErrorMessage      string           `json:"errorMessage,omitempty"`
	HaltReason        string           `json:"haltReason,omitempty"`
	HaltExplination   string           `json:"haltExplination,omitempty"`
	ValidationMessage string           `json:"validationMessage,omitempty"`
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
