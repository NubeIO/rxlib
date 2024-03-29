package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/plugins"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/NubeIO/scheduler"
	"sync"
)

type Runtime interface {
	Get() []Object
	GetObjectsConfig() []*runtime.ObjectConfig
	GetObjectConfig(uuid string) *runtime.ObjectConfig

	GetObjectValues(objectUUID string) []*runtime.PortValue
	GetObjectsValues() map[string][]*runtime.PortValue

	Delete() string
	GetByUUID(uuid string) Object
	GetFirstByID(objectID string) Object
	GetAllByID(objectID string) []Object
	GetFirstByName(name string) Object
	GetAllByName(name string) []Object

	GetChildObjectsByWorkingGroup(objectUUID, workingGroup string) []Object

	GetChildObjects(parentUUID string) []Object
	GetAllObjectValues() []*ObjectValue
	AddObject(object Object)
	Command(cmd *ExtendedCommand) *runtime.CommandResponse
	CommandObject(cmd *ExtendedCommand) *CommandResponse

	GetTreeMapRoot() *runtime.ObjectsRootMap
	GetAncestorTreeByUUID(objectUUID string) *AncestorTreeNode
	GetChilds(objectUUID string) *AncestorTreeNode

	AllPlugins() []*plugins.Export

	GetObjectsPallet() *PalletTree

	Scheduler() scheduler.Scheduler
}

type RuntimeOpts struct {
	Scheduler scheduler.Scheduler
}

func NewRuntime(objs []Object, opts *RuntimeOpts) Runtime {
	r := &RuntimeImpl{
		tree: &tree{},
	}
	r.objects = objs
	r.scheduler = opts.Scheduler
	if r.objects == nil {
		panic("NewRuntime []Object can not be empty")
	}
	return r
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

func (inst *RuntimeImpl) GetObjectConfig(uuid string) *runtime.ObjectConfig {
	return serializeCurrentFlowArray(inst.GetByUUID(uuid))
}

func NewCommandResponse() *runtime.CommandResponse {
	return &runtime.CommandResponse{}
}

type CommandResponse struct {
	SenderID         string                  `json:"senderID,omitempty"` // if sent from another ROS instance
	Count            int                     `json:"count,omitempty"`
	Objects          []Object                `json:"-,omitempty"`
	SerializeObjects []*runtime.ObjectConfig `json:"serializeObjects,omitempty"`
	MapPorts         map[string][]*Port      `json:"mapPorts,omitempty"`
	MapStrings       map[string]string       `json:"mapStrings,omitempty"`
	Float            float64                 `json:"number,omitempty"`
	Bool             bool                    `json:"boolean,omitempty"`
	Error            string                  `json:"error,omitempty"`
	ReturnType       string                  `json:"returnType,omitempty"`
	Any              []byte                  `json:"any,omitempty"`
	CommandResponse  []*CommandResponse      `json:"response,omitempty"`
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

func (inst *RuntimeImpl) Get() []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst.objects
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

// let obj = RQL.AllObjects().Where("histories").Name("uuid").Condition("==", "hist_history").First()
// let obj = RQL.AllObjects().Where("objects").Name("name").Condition("==", "abc").SerialObjects()
func (inst *RuntimeImpl) histories() *RuntimeImpl {
	//var filtered = make(map[string]Object)
	//for _, obj := range inst.objectsFiltered {
	//	extension := obj.GetRequiredExtensionByName("history")
	//	if extension != nil {
	//		filtered[obj.GetUUID()] = obj
	//	}
	//}
	//inst.objectsFiltered = filtered
	return inst
}

const (
	operatorEqual    = "=="
	operatorNotEqual = "!="
)

var operatorValues = []string{operatorEqual, operatorNotEqual}

const (
	fieldName = "name"
	fieldUUID = "uuid"
)

var fieldValues = []string{fieldName, fieldUUID}

func compareObject(object Object, field, operator, value string) bool {
	var fieldValue string
	switch field {
	case "name":
		fieldValue = object.GetName()
	case "uuid":
		fieldValue = object.GetUUID()
	case "id":
		fieldValue = object.GetID()
	case "objectID":
		fieldValue = object.GetID()
	}
	switch operator {
	case "==":
		return fieldValue == value
	case "!=":
		return fieldValue != value
	}

	return false
}

func compareHist(object Object, field, operator, value string) bool {
	//Obj.GetHistoryManager().AllHistoriesByObjectUUID()
	//switch operator {
	//case "==":
	//	return fieldValue == value
	//case "!=":
	//	return fieldValue != value
	//}
	//return false
	return false
}

func comparePorts(port *Port, field, operator, value string) bool {
	var fieldValue string

	switch field {
	case "name":
		fieldValue = port.GetName()
	case "uuid":
		fieldValue = port.GetUUID()

	}

	switch operator {
	case "==":
		return fieldValue == value
	case "!=":
		return fieldValue != value
	}

	return false
}
