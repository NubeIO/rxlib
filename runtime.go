package rxlib

import (
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib/libs/history"
	"github.com/NubeIO/rxlib/libs/pglib"
	"github.com/NubeIO/rxlib/libs/restc"
	systeminfo "github.com/NubeIO/rxlib/libs/system"
	"github.com/NubeIO/rxlib/plugins"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/NubeIO/scheduler"
	"sync"
)

type Runtime interface {
	// Get all objects []Object
	Get() []Object
	// AddObjects adds objects to runtime
	AddObjects([]Object)
	// ToObjectsConfig converts to ObjectConfig, used when needed as JSON
	ToObjectsConfig(objects []Object) []*runtime.ObjectConfig
	// GetObjectsUUIDs returns UUIDs of objects
	GetObjectsUUIDs(objects []Object) []string
	// GetObjectsConfig returns config of all objects
	GetObjectsConfig() []*runtime.ObjectConfig
	// GetObjectConfig returns config of an object by UUID
	GetObjectConfig(uuid string) *runtime.ObjectConfig
	// GetObjectConfigByID returns config of an object by ID
	GetObjectConfigByID(objectID string) *runtime.ObjectConfig

	// GetObjectValues returns values of an object by UUID
	GetObjectValues(objectUUID string) []*runtime.PortValue
	// GetObjectsValues returns values of all objects
	GetObjectsValues() map[string][]*runtime.PortValue
	// GetObjectsValuesPaginate returns paginated values of objects
	GetObjectsValuesPaginate(parentUUID string, pageNumber, pageSize int) *ObjectValuesPagination

	// ObjectsPagination paginates objects
	ObjectsPagination(pageNumber, pageSize int) *ObjectPagination
	// PaginateGetAllByID paginates objects by ID
	PaginateGetAllByID(objectID string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetChildObjects paginates child objects
	PaginateGetChildObjects(parentUUID string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetAllByName paginates objects by name
	PaginateGetAllByName(name string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetChildObjectsByWorkingGroup paginates child objects by working group
	PaginateGetChildObjectsByWorkingGroup(objectUUID, workingGroup string, pageNumber, pageSize int) *ObjectPagination

	// Delete deletes runtime
	Delete() string
	// GetByUUID gets object by UUID; eg GetByUUID("abc").GetName(), GetByUUID("abc").GetInputs()
	GetByUUID(uuid string) Object
	// GetFirstByID gets first object by ID eg GetFirstByID("abc").GetID(), GetFirstByID("abc").GetTags()
	GetFirstByID(objectID string) Object
	// GetAllByID gets all objects by ID
	GetAllByID(objectID string) []Object
	// GetFirstByName gets first object by name
	GetFirstByName(name string) Object
	// GetAllByName gets all objects by name
	GetAllByName(name string) []Object

	// GetChildObjectsByWorkingGroup gets child objects by working group
	GetChildObjectsByWorkingGroup(objectUUID, workingGroup string) []Object
	// GetChildObjects gets child objects
	GetChildObjects(parentUUID string) []Object
	// GetAllObjectValues gets all object values
	GetAllObjectValues() []*ObjectValue
	// AddObject adds an object
	AddObject(object Object)
	// Command executes a command
	Command(cmd *ExtendedCommand) *runtime.CommandResponse
	// CommandObject executes a command for an object
	CommandObject(cmd *ExtendedCommand) *CommandResponse

	// GetTreeMapRoot gets the root of the object tree map
	GetTreeMapRoot() *runtime.ObjectsRootMap
	// GetAncestorTreeByUUID gets ancestor tree by UUID
	GetAncestorTreeByUUID(objectUUID string) *AncestorTreeNode
	// GetChilds gets child nodes of an object
	GetChilds(objectUUID string) *AncestorTreeNode

	// AllPlugins returns all plugins
	AllPlugins() []*plugins.Export

	// GetObjectsPallet returns objects pallet tree
	GetObjectsPallet() *PalletTree

	// Scheduler gets scheduler
	Scheduler() scheduler.Scheduler

	// ExprWithError run a system query and returns an error. eg; Expr("filter(objects, .GetID() == "rubix-manager""))  see docs https://github.com/expr-lang/expr
	ExprWithError(query string) (any, error)

	// Expr run a system query. eg; Expr("filter(objects, .GetID() == "rubix-manager""))  see docs https://github.com/expr-lang/expr
	Expr(query string) any

	// System get host info, networking, memory and stats. eg; System().GetIP()
	System() systeminfo.System

	// SystemInfo return the host machine info eg; update and timezone
	SystemInfo() *systeminfo.Info

	// HistoryManager get ros history manager. eg; HistoryManager().AllHistories()
	HistoryManager() history.Manager

	// ToStringArray Conversions
	ToStringArray(interfaces interface{}) []string

	// Rest is a rest client for making HTTP requests, Example Rest().Execute("GET", "http://localhost:8080/api")
	Rest() restc.Rest

	// Publish a mqtt message
	Publish(topic string, body interface{}) (err string)

	// Iam is used for discovery ROS instances that are connected to command MQTT broker. eg; {"key": "command", "args": ["run", "whois"], "data": {"start": "1", "finish": "200", "global": "true"}}
	Iam(rangeStart, finish int) Object

	// DB access the postgres database; eg DG().Select("select *")
	DB() pglib.PG

	// ObjectSync sync all the objects to the postgres db;
	ObjectSync(forceSync bool, opts *SyncOptions) error
	// HistorySync sync all the histories to the postgres db;
	HistorySync(forceSync bool, opts *SyncOptions) error
}

type RuntimeOpts struct {
	Scheduler  scheduler.Scheduler
	MQTTClient mqttwrapper.MQTT
}

func NewRuntime(objs []Object, opts *RuntimeOpts) Runtime {
	r := &RuntimeImpl{
		tree:       &tree{},
		mqttClient: opts.MQTTClient,
	}
	r.objects = objs
	r.scheduler = opts.Scheduler
	r.hist = history.NewHistoryManager("ros")
	connString := "postgresql://postgres:postgres@localhost/postgres"
	db, err := pglib.New(connString)
	if err != nil {
		fmt.Println("runtime init DB: ", err)
	}
	r.db = db
	r.rest = restc.New()
	return r
}

func (inst *RuntimeImpl) Get() []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst.objects
}

type RuntimeImpl struct {
	objects         []Object
	objectsFiltered []Object
	PluginsExport   []*plugins.Export
	err             error // To handle errors in query chain
	where           string
	field           string
	mutex           sync.RWMutex
	response        *CommandResponse
	parsedCommand   *ParsedCommand
	command         *ExtendedCommand
	tree            *tree
	addedObject     bool
	scheduler       scheduler.Scheduler
	hist            history.Manager
	db              pglib.PG
	rest            restc.Rest
	mqttClient      mqttwrapper.MQTT
}

func (inst *RuntimeImpl) Publish(topic string, body interface{}) string {
	if inst.mqttClient == nil {
		return "client is empty"
	}
	err := inst.mqttClient.Publish(topic, body)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (inst *RuntimeImpl) AddObjects(objects []Object) {
	inst.objects = objects
}

func (inst *RuntimeImpl) HistoryManager() history.Manager {
	return inst.hist
}

func (inst *RuntimeImpl) DB() pglib.PG {
	return inst.db
}

func (inst *RuntimeImpl) Rest() restc.Rest {
	return inst.rest
}

func (inst *RuntimeImpl) System() systeminfo.System {
	return systeminfo.NewSystem()
}

func (inst *RuntimeImpl) SystemInfo() *systeminfo.Info {
	return inst.System().Info()
}

func (inst *RuntimeImpl) Scheduler() scheduler.Scheduler {
	return inst.scheduler
}

// GetChildObjectsByWorkingGroup
// for example get all the childs object for working group "rubix"
func (inst *RuntimeImpl) GetChildObjectsByWorkingGroup(objectUUID, workingGroup string) []Object {
	var out []Object
	for _, object := range inst.Get() {
		if object.GetUUID() == objectUUID {
			if object.GetWorkingGroup() == workingGroup {
				out = append(out, object)
			}
		}
	}
	return out
}

func (inst *RuntimeImpl) GetObjectValues(objectUUID string) []*runtime.PortValue {
	obj := inst.GetByUUID(objectUUID)
	if obj == nil {
		return nil
	}
	var out []*runtime.PortValue
	inputs := obj.GetInputs()
	for _, port := range inputs {
		out = append(out, obj.GetPortValue(port.GetID()))
	}
	outputs := obj.GetOutputs()
	for _, port := range outputs {
		out = append(out, obj.GetPortValue(port.GetID()))
	}
	return out
}

func (inst *RuntimeImpl) GetObjectsValues() map[string][]*runtime.PortValue {
	out := make(map[string][]*runtime.PortValue)
	for _, object := range inst.Get() {
		out[object.GetUUID()] = inst.GetObjectValues(object.GetUUID())
	}
	return out
}

func (inst *RuntimeImpl) GetObjectsConfig() []*runtime.ObjectConfig {
	return SerializeCurrentFlowArray(inst.Get())
}

func (inst *RuntimeImpl) ToObjectsConfig(objects []Object) []*runtime.ObjectConfig {
	return SerializeCurrentFlowArray(objects)
}

func (inst *RuntimeImpl) GetObjectsUUIDs(objects []Object) []string {
	var out []string
	for _, object := range objects {
		out = append(out, object.GetUUID())
	}
	return out
}

func (inst *RuntimeImpl) GetObjectConfig(uuid string) *runtime.ObjectConfig {
	return serializeCurrentFlowArray(inst.GetByUUID(uuid))
}

func (inst *RuntimeImpl) GetObjectConfigByID(objectID string) *runtime.ObjectConfig {
	object := inst.GetFirstByID(objectID)
	if object == nil {
		return nil
	}
	return serializeCurrentFlowArray(object)
}

func NewCommandResponse() *runtime.CommandResponse {
	return &runtime.CommandResponse{}
}

type CommandResponse struct {
	SenderID         string                    `json:"senderID,omitempty"` // if sent from another ROS instance
	Count            int                       `json:"count,omitempty"`
	Objects          []Object                  `json:"-,omitempty"`
	SerializeObjects []*runtime.ObjectConfig   `json:"serializeObjects,omitempty"`
	MapPorts         map[string][]*Port        `json:"mapPorts,omitempty"`
	MapStrings       map[string]string         `json:"mapStrings,omitempty"`
	Float            float64                   `json:"number,omitempty"`
	Bool             bool                      `json:"boolean,omitempty"`
	Error            string                    `json:"error,omitempty"`
	ReturnType       string                    `json:"returnType,omitempty"`
	Byte             []byte                    `json:"byte,omitempty"`
	CommandResponse  []*CommandResponse        `json:"response,omitempty"`
	ObjectPagination *runtime.ObjectPagination `json:"objectPagination,omitempty"`
	ObjectTree       *runtime.ObjectsRootMap   `json:"objectTree,omitempty"`
	Data             any                       `json:"data"`
}

func (inst *RuntimeImpl) GetTreeMapRoot() *runtime.ObjectsRootMap {
	if !inst.addedObject {
		inst.tree.addObjects(inst.objects)
	}
	inst.addedObject = true
	return inst.tree.GetTreeMapRoot()
}

func (inst *RuntimeImpl) GetAncestorTreeByUUID(objectUUID string) *AncestorTreeNode {
	if !inst.addedObject {
		inst.tree.addObjects(inst.objects)
	}
	inst.addedObject = true
	return inst.tree.GetAncestorTreeByUUID(objectUUID)
}

func (inst *RuntimeImpl) GetChilds(objectUUID string) *AncestorTreeNode {
	return inst.tree.GetChilds(objectUUID)
}

func (inst *RuntimeImpl) Delete() string {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	c := len(inst.objects)
	inst.objects = nil
	d := len(inst.objects)
	return fmt.Sprintf("count deleted: %d current: %d", c, d)
}

func (inst *RuntimeImpl) AddObject(object Object) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.objects = append(inst.objects, object)
}

func (inst *RuntimeImpl) GetAllByID(objectID string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetID() == objectID {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetByUUID(uuid string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, object := range inst.objects {
		if object.GetUUID() == uuid {
			return object
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetAllByName(name string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetName() == name {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetChildObjects(parentUUID string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetParentUUID() == parentUUID {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetFirstByID(objectID string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, obj := range inst.objects {
		if obj.GetID() == objectID {
			return obj
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetFirstByName(name string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, obj := range inst.objects {
		if obj.GetName() == name {
			return obj
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetAllObjectValues() []*ObjectValue {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	allObjects := inst.Get()
	nodeValues := make([]*ObjectValue, len(allObjects))
	for _, node := range allObjects {
		nv := node.GetAllPorts()
		if nv == nil {
			continue
		}
		portValue := &ObjectValue{
			ObjectId:   node.GetID(),
			ObjectUUID: node.GetUUID(),
			Ports:      nv,
		}
		nodeValues = append(nodeValues, portValue)
	}
	return nodeValues
}

func (inst *RuntimeImpl) objectsFilteredIsNil() bool {
	if inst.objectsFiltered == nil {
		return true
	}
	return false
}
