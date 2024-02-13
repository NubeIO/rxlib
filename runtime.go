package rxlib

import (
	"fmt"
	"strings"
	"sync"
)

type ObjectResponse struct {
	Obj Object
	Err error
}

func (inst *ObjectResponse) Object() Object {
	return inst.Obj
}

func (inst *ObjectResponse) Error() string {
	if inst.Err != nil {
		return inst.Err.Error()
	}
	return ""
}

func (inst *ObjectResponse) GetError() error {
	return inst.Err
}

type Runtime interface {
	Get() []Object
	GetByUUID(uuid string) Object
	GetFirstByID(objectID string) Object
	GetFirstByName(name string) Object

	GetForeignObject(objectUUID string) Object
	CheckForeignObjectOutputExists(objectUUID, portID string) *Port
	CheckForeignObjectInputExists(objectUUID, portID string) *Port
	GetObjectsByType(objectID string) []Object // for example get all math/add Object

	GetAllObjectValues() []*ObjectValue

	AddObject(object Object)
	AddObjectResp(object Object) *ObjectResponse

	Where(attribute string) *RuntimeImpl
	Field(field string) *RuntimeImpl
	Condition(operator, value string) *RuntimeImpl
	First() Object
	Filtered() []Object

	Query(query string, r Runtime) []Object
}

func NewRuntime(objs []Object) Runtime {
	r := &RuntimeImpl{}
	r.objects = objs
	if r.objects == nil {
		panic("NewRuntime Obj map can not be empty")
	}
	return r
}

type RuntimeImpl struct {
	objects         []Object
	objectsFiltered []Object
	err             error // To handle errors in query chain
	where           string
	field           string
	mutex           sync.RWMutex
}

func (inst *RuntimeImpl) Query(query string, r Runtime) []Object {
	var allResults []Object
	orSegments := strings.Split(query, "OR") // Split the query into OR segments

	for _, orSegment := range orSegments {
		orSegment = strings.TrimSpace(strings.Trim(orSegment, "()"))
		andConditions := strings.Split(orSegment, "AND") // Split the OR segment into AND conditions

		var segmentResults []Object // Temporarily store results for this OR segment
		if len(andConditions) == 0 {
			continue
		}

		segmentResults = r.Get() // Start with all objects for the first condition
		for _, andCondition := range andConditions {
			andCondition = strings.TrimSpace(andCondition)
			// Identify the operator and split the condition
			var operator string
			if strings.Contains(andCondition, "==") {
				operator = "=="
			} else if strings.Contains(andCondition, "!=") {
				operator = "!="
			}
			parts := strings.SplitN(andCondition, operator, 2)
			if len(parts) != 2 {
				fmt.Println("Invalid condition:", andCondition)
				continue
			}

			field := strings.TrimSpace(parts[0])
			// Remove the 'objects:' prefix and any extraneous characters like parentheses
			field = strings.ReplaceAll(field, "objects:", "")
			field = strings.Trim(field, "() ")

			value := strings.TrimSpace(parts[1])
			// Here's the key adjustment: ensure we also trim the closing parenthesis from the value
			value = strings.Trim(value, "() ")

			// Filter objects that match the current condition
			var matchesForThisCondition []Object
			for _, obj := range segmentResults {
				if matchesCondition(obj, field, operator, value) {
					matchesForThisCondition = append(matchesForThisCondition, obj)
				}
			}

			// Update segmentResults to narrow down the matches for the next AND condition
			segmentResults = matchesForThisCondition
		}

		// Add unique matches from this segment to allResults
		for _, match := range segmentResults {
			if !containsObject(allResults, match) {
				allResults = append(allResults, match)
			}
		}
	}

	return allResults
}

// matchesCondition checks if an object matches the given field and value.
func matchesCondition(obj Object, field, operator, value string) bool {
	var fieldValue string
	value = strings.Trim(value, " )")
	switch field {
	case "uuid":
		fieldValue = obj.GetUUID()
	case "category":
		fieldValue = obj.GetCategory()
	case "id":
		fieldValue = obj.GetID()
	case "name":
		fieldValue = obj.GetName()
	}
	switch operator {
	case "==":
		return fieldValue == value
	case "!=":
		return fieldValue != value
	}

	return false
}

// containsObject checks if the object is already in the slice.
func containsObject(slice []Object, obj Object) bool {
	for _, item := range slice {
		if item.GetUUID() == obj.GetUUID() {
			return true
		}
	}
	return false
}

func (inst *RuntimeImpl) AddObjectResp(object Object) *ObjectResponse {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.objects = append(inst.objects, object)
	return nil
}

func (inst *RuntimeImpl) AddObject(object Object) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.objects = append(inst.objects, object)
}

func (inst *RuntimeImpl) GetObjectsByType(objectID string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out = make(map[string]Object)
	for _, obj := range inst.objects {
		if obj.GetID() == objectID {
			out[obj.GetUUID()] = obj
		}
	}
	return nil
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

func (inst *RuntimeImpl) GetForeignObject(objectUUID string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	obj := inst.GetByUUID(objectUUID)
	return obj
}

func (inst *RuntimeImpl) CheckForeignObjectInputExists(objectUUID, portID string) *Port {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	obj := inst.GetByUUID(objectUUID)
	for _, port := range obj.GetInputs() {
		if port.GetID() == portID {
			return port
		}
	}
	return nil
}

func (inst *RuntimeImpl) CheckForeignObjectOutputExists(objectUUID, portID string) *Port {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	obj := inst.GetByUUID(objectUUID)
	for _, port := range obj.GetOutputs() {
		if port.GetID() == portID {
			return port
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

// QUERY

func (inst *RuntimeImpl) Where(attribute string) *RuntimeImpl {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	if inst.err != nil {
		return inst // Skip processing if there's an Err
	}
	inst.where = attribute

	var filtered []Object
	for _, obj := range inst.objects {
		switch attribute {
		case "histories":
			inst.histories()
			filtered = inst.objects
		case "objects":
			filtered = inst.objects
		case "inputs":
			if len(obj.GetInputs()) > 0 {
				filtered = append(filtered, obj)
			}
		case "outputs":
			if len(obj.GetOutputs()) > 0 {
				filtered = append(filtered, obj)
			}
		}
	}
	inst.objectsFiltered = filtered
	return inst
}

func (inst *RuntimeImpl) Field(field string) *RuntimeImpl {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.field = field
	return inst
}

func (inst *RuntimeImpl) Condition(operator, value string) *RuntimeImpl {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	if inst.err != nil {
		return inst // Skip processing if there's an Err
	}
	var filtered []Object
	for _, obj := range inst.objectsFiltered {
		switch inst.where {
		case "histories":
			// Comparing Obj fields
			if compareHist(obj, inst.field, operator, value) {

				filtered = append(filtered, obj)
			}
		case "objects":
			// Comparing Obj fields
			//fmt.Println(compareObject(obj, inst.field, operator, value), inst.field, value, obj.GetID())
			if compareObject(obj, inst.field, operator, value) {
				filtered = append(filtered, obj)
			}

		case "inputs":
			// Comparing fields of input ports
			for _, port := range obj.GetInputs() {
				if comparePorts(port, inst.field, operator, value) {
					filtered = append(filtered, obj)
					break // Found a matching input, no need to check further
				}
			}

		case "outputs":
			// Comparing fields of output ports
			for _, port := range obj.GetOutputs() {
				if comparePorts(port, inst.field, operator, value) {
					filtered = append(filtered, obj)
					break // Found a matching output, no need to check further
				}
			}
		}
	}

	inst.objectsFiltered = filtered
	return inst
}

func (inst *RuntimeImpl) First() Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	if inst.objectsFilteredIsNil() {
		return nil
	}
	for _, value := range inst.objectsFiltered {
		return value
	}
	return nil
}

func (inst *RuntimeImpl) Filtered() []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst.objectsFiltered
}

func (inst *RuntimeImpl) objectsFilteredIsNil() bool {
	if inst.objectsFiltered == nil {
		return true
	}
	return false
}

// let obj = RQL.AllObjects().Where("histories").Field("uuid").Condition("==", "hist_history").First()
// let obj = RQL.AllObjects().Where("objects").Field("name").Condition("==", "abc").SerialObjects()
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

//-------------CONNECTIONS------------------

func (inst *RuntimeImpl) AddConnection(sourceUUID, sourcePort, targetUUID, targetPort string) Object {
	//connection, c, err := NewConnection(sourceUUID, sourcePort, targetUUID, targetPort)
	//if err != nil {
	//	return nil
	//}
	return nil
}
