package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/convert"
	"github.com/NubeIO/rxlib/priority"
	"strconv"
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
	Delete() string
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

	Query(query string) []Object

	CommandObject(cmd *Command) any
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

// --------------- COMMANDS ----------------

func (inst *RuntimeImpl) CommandObject(cmd *Command) any {
	if cmd == nil {
		return nil
	}
	commandType := cmd.CommandType
	field := cmd.Field
	uuidName := cmd.FieldEntry

	var result any
	if cmd.Query != "" {
		fmt.Println(cmd.Query)
		objects := inst.Query(cmd.Query)
		fmt.Printf("found objects: %d from query -type: %s -cmd: %s \n", len(objects), commandType, cmd.CommandName)
		switch strings.ToLower(string(commandType)) {
		case "get":
			result = handleGetCommandForMultipleObjects(cmd, objects)
		case "set":
			result = handleSetCommandForMultipleObjects(cmd, objects)
		default:
			result = fmt.Errorf("unknown command type: %s", commandType)
		}
	} else {
		var object Object
		if field == "uuid" {
			object = inst.GetByUUID(uuidName)
		} else if field == "name" {
			object = inst.GetFirstByName(uuidName)
		}

		if object == nil {
			return fmt.Errorf("object not found field: %s fieldEntry: %s", field, uuidName)
		}

		switch strings.ToLower(string(commandType)) {
		case "get":
			result = handleGetCommand(cmd, object)
		case "set":
			result = handleSetCommand(cmd, object)
		default:
			result = fmt.Errorf("unknown command type: %s", commandType)
		}
	}

	return result
}

func handleGetCommandForMultipleObjects(cmd *Command, objects []Object) any {
	var results []any
	for _, object := range objects {
		result := handleGetCommand(cmd, object)
		results = append(results, result)
	}
	return results
}

func handleSetCommandForMultipleObjects(cmd *Command, objects []Object) any {
	var results []any
	for _, object := range objects {
		result := handleSetCommand(cmd, object)
		results = append(results, result)
	}
	return results
}

func handleGetCommand(cmd *Command, object Object) any {
	fmt.Println(cmd.CommandName)
	switch strings.ToLower(cmd.CommandName) {
	case "object":
		return object

	case "uuid":
		return object.GetUUID()
	case "name":
		return object.GetName()
	case "inputs":
		return object.GetInputs()
	case "input":
		return handlePort(cmd.Args, object)
	default:
		return fmt.Errorf("unknown get command: %s", cmd.CommandName)
	}

}

func handlePort(args []string, object Object) any {
	fmt.Println("&&&&&&&&&&&&&&&&")
	if len(args) < 1 {
		return fmt.Errorf("getInput command requires an argument")
	}
	fmt.Println(args)
	if len(args) == 3 {
		arg1 := args[0] // id
		arg2 := args[1] // value
		arg3 := args[2] // data, value, name
		if arg1 == "id" {
			var port *Port
			port = object.GetInput(arg2)
			if arg3 == "name" {
				return port.GetName()
			}
			if arg3 == "value" || arg3 == "data" {
				return port.GetValueDisplay()
			}
		}
	}
	return fmt.Errorf("unknown get command")
}

func handleSetCommand(cmd *Command, object Object) any {
	switch strings.ToLower(cmd.CommandName) {
	case "name":
		if len(cmd.Args) < 1 {
			return fmt.Errorf("setname command requires an argument")
		}
		return object.SetName(cmd.Args[0])
	case "input":
		if len(cmd.Args) < 1 {
			return fmt.Errorf("getInput command requires an argument")
		}
		if len(cmd.Args) == 4 {
			arg1 := cmd.Args[0] // id
			arg2 := cmd.Args[1] // value
			arg3 := cmd.Args[2] // data, value, name
			arg4 := cmd.Args[3] // 22.5, "new name"

			if arg1 == "id" {
				var port *Port
				port = object.GetInput(arg2)
				if arg3 == "name" {
					return port.SetName(arg4)
				}
				if arg3 == "write" {
					if port.GetDataType() == priority.TypeFloat {
						f := convert.AnyToFloatPointer(arg4)
						if f == nil {
							return "was unable to parse value as type float"
						}
						err := port.Write(f)
						if err != nil {
							return err.Error()
						}
						return fmt.Sprintf("object: %s updated ok port: %s value: %s", object.GetName(), arg2, arg4)
					}
					return port.Write(arg4)
				}
			}
		}
	default:
		return fmt.Errorf("unknown set command: %s", cmd.CommandName)
	}
	return fmt.Errorf("unknown get command")

}

func (inst *RuntimeImpl) Delete() string {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	c := len(inst.objects)
	inst.objects = nil
	d := len(inst.objects)
	return fmt.Sprintf("count deleted: %d current: %d", c, d)
}

func (inst *RuntimeImpl) Query(query string) []Object {
	var allResults []Object
	query, limit := extractAndRemoveLimit(query) // get the query limit
	orSegments := strings.Split(query, "OR")

	for _, orSegment := range orSegments {
		orSegment = strings.TrimSpace(strings.Trim(orSegment, "()"))
		andConditions := strings.Split(orSegment, "AND") // Split the OR segment into AND conditions

		var segmentResults []Object // Temporarily store results for this OR segment
		if len(andConditions) == 0 {
			continue
		}
		segmentResults = inst.Get() // Start with all objects for the first condition

		segmentResults = filterObjectsByCondition(segmentResults, andConditions)
		allResults = addUniqueMatches(allResults, segmentResults)

	}

	if limit >= 0 && len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults
}

func addUniqueMatches(allResults []Object, segmentResults []Object) []Object {
	for _, match := range segmentResults {
		if !containsObject(allResults, match) {
			allResults = append(allResults, match)
		}
	}
	return allResults
}

func filterObjectsByCondition(segmentResults []Object, andConditions []string) []Object {
	var filteredResults []Object

	for _, andCondition := range andConditions {
		field, operator, value := extractFieldAndValue(andCondition)
		if field == "" {
			fmt.Println("Invalid condition:", andCondition)
			continue
		}

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

	filteredResults = segmentResults
	return filteredResults
}

func matchesCondition(obj Object, field, operator, value string) bool {
	switch {
	case strings.HasPrefix(field, "inputs:"):
		return matchesPortCondition(obj.GetInputs(), field, operator, value)
	case strings.HasPrefix(field, "outputs:"):
		return matchesPortCondition(obj.GetOutputs(), field, operator, value)
	default:
		return matchesObjectCondition(obj, field, operator, value)
	}
}

func extractFieldAndValue(condition string) (string, string, string) {
	condition = strings.TrimSpace(condition)
	var operator string
	if strings.Contains(condition, "==") {
		operator = "=="
	} else if strings.Contains(condition, "!=") {
		operator = "!="
	} else if strings.Contains(condition, ">=") {
		operator = ">="
	} else if strings.Contains(condition, ">") {
		operator = ">"
	}

	parts := strings.SplitN(condition, operator, 2)
	if len(parts) != 2 {
		return "", "", ""
	}

	field := strings.TrimSpace(parts[0])
	// remove the 'objects:' prefix and any extraneous characters like parentheses
	field = strings.ReplaceAll(field, "objects:", "")
	field = strings.Trim(field, "() ")

	value := strings.TrimSpace(parts[1])
	// ensure we trim the closing parenthesis from the value
	value = strings.Trim(value, "() ")

	return field, operator, value
}

func matchesPortCondition(ports []*Port, field, operator, value string) bool {
	// extract the attribute we're interested in (e.g., "name", "id", "uuid", "value")
	parts := strings.SplitN(field, ":", 2)
	if len(parts) != 2 {
		return false
	}
	fieldAttribute := parts[1]
	for _, port := range ports {
		if matchesPortAttributeCondition(port, fieldAttribute, operator, value) {
			return true
		}
	}
	return false
}

func matchesPortAttributeCondition(port *Port, attribute, operator, value string) bool {
	switch attribute {
	case "value":
		if value == "null" { // get all ports that are null/nil
			isNull := port.GetValue().IsNull()
			if isNull {
				return true
			}
			return false
		}
		compareValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false // Invalid comparison value, return false
		}
		portValue := port.GetValue().GetFloat()
		return compareFloat64(portValue, operator, compareValue)
	default:
		return false // Invalid attribute, return false
	}
}

func compareFloat64(portValue float64, operator string, compareValue float64) bool {
	switch operator {
	case "==":
		return portValue == compareValue
	case "!=":
		return portValue != compareValue
	case ">":
		return portValue > compareValue
	case ">=":
		return portValue >= compareValue
	case "<":
		return portValue < compareValue
	case "<=":
		return portValue <= compareValue
	default:
		return false // Invalid operator, return false
	}
}

func matchesObjectCondition(obj Object, field, operator, value string) bool {
	// handle other fields (uuid, category, id, name)
	var fieldValue string
	switch field {
	case "uuid":
		fieldValue = obj.GetUUID()
	case "category":
		fieldValue = obj.GetCategory()
	case "id":
		fieldValue = obj.GetID()
	case "name":
		fieldValue = obj.GetName()
	default:
		return false
	}

	// Evaluate the condition for non-port fields
	switch operator {
	case "==":
		return fieldValue == value
	case "!=":
		return fieldValue != value
	default:
		return false
	}
}

func extractAndRemoveLimit(query string) (string, int) {
	limit := -1
	if strings.Contains(query, "limit:") {
		parts := strings.Split(query, "limit:")
		if len(parts) == 2 {
			limitString := strings.TrimSpace(parts[1])
			query = strings.Replace(query, "limit:"+limitString, "", 1) // remove the limit substring from the query
			limit, _ = strconv.Atoi(limitString)
		}
	}
	return query, limit
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

//-------------CONNECTIONS------------------

func (inst *RuntimeImpl) AddConnection(sourceUUID, sourcePort, targetUUID, targetPort string) Object {
	//connection, c, err := NewConnection(sourceUUID, sourcePort, targetUUID, targetPort)
	//if err != nil {
	//	return nil
	//}
	return nil
}
