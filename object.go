package rxlib

import (
	"github.com/NubeIO/schema"
	"github.com/gin-gonic/gin"
)

type Object interface {
	New(objectUUID string, object Object, settings *Settings, opts ...any) Object
	Start()
	Delete()
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool

	// object details/info
	GetID() string
	GetUUID() string
	GetParentUUID() string

	// info
	ObjectInfo
	SetInfo(info *Info)
	GetInfo() *Info

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
	HaltExplanation   string           `json:"haltExplanation,omitempty"`
	ValidationMessage string           `json:"validationMessage,omitempty"`
}

type ObjectValue struct {
	ObjectId   string  `json:"objectId"`
	ObjectUUID string  `json:"objectUUID"`
	Ports      []*Port `json:"ports"`
}

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	ObjectName string   `json:"objectName"` // comes from UI need to set in objectInfo
	ParentUUID string   `json:"parentUUID"` // comes from UI need to set in objectInfo
	Position   Position `json:"position"`
}
