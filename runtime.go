package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
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

	Where(attribute string) *RuntimeImpl
	Field(field string) *RuntimeImpl
	Condition(operator, value string) *RuntimeImpl
	First() Object
	Filtered() []Object

	Query(query string) []Object

	CommandObject(cmd *Command) *CommandResponse
}

func NewRuntime(objs []Object) Runtime {
	r := &RuntimeImpl{}
	r.objects = objs
	if r.objects == nil {
		panic("NewRuntime []Object can not be empty")
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
	cmd             *CommandResponse
	parsedCommand   *parsedCommand
}

type CommandResponse struct {
	Object           Object             `json:"object,omitempty"`
	SerializeObject  *ObjectConfig      `json:"serializeObject"`
	Objects          []Object           `json:"objects,omitempty"`
	SerializeObjects []*ObjectConfig    `json:"serializeObjects"`
	Ports            []*Port            `json:"ports,omitempty"`
	MapPorts         map[string][]*Port `json:"mapPorts,omitempty"`
	String           string             `json:"string,omitempty"`
	Float            *float64           `json:"number,omitempty"`
	Bool             *bool              `json:"boolean,omitempty"`
	Error            error              `json:"error,omitempty"`
	ReturnType       string             `json:"returnType,omitempty"`
}

func (inst *RuntimeImpl) CommandObject(command *Command) *CommandResponse {
	inst.cmd = &CommandResponse{}
	if command == nil {
		inst.cmd.Error = fmt.Errorf("command cannot be nil")
		return inst.cmd
	}
	returnType, parsedArgs := commandReturnType(command)
	inst.cmd.ReturnType = returnType
	var returnAsJSON bool
	if parsedArgs.ReturnAs == commandJSON {
		returnAsJSON = true
	}
	if parsedArgs.ReturnAs != "" {
		fmt.Printf("type: %s thing: %s type-return: %s --as: %s \n", command.CommandType, command.Thing, returnType, parsedArgs.ReturnAs)
	} else {
		fmt.Printf("type: %s thing: %s type-return: %s \n", command.CommandType, command.Thing, returnType)
	}

	inst.parsedCommand = parsedArgs
	if command.Query != "" {
		objects := inst.Query(command.Query)
		if returnAsJSON {
			inst.cmd.SerializeObjects = SerializeCurrentFlowArray(objects)
		} else {
			inst.cmd.Objects = objects
		}
		return inst.cmd
	}
	uuid := parsedArgs.UUID
	name := parsedArgs.Name
	byType := parsedArgs.Type
	var object Object
	if uuid != "" {
		object = inst.GetByUUID(uuid)
	} else if name != "" {
		object = inst.GetFirstByName(name)
	} else if byType != "" {
		object = inst.GetFirstByID(name)
	}
	if object == nil {
		inst.cmd.Error = fmt.Errorf("object not found by uuid: %s, name: %s", uuid, name)
		return inst.cmd
	}

	switch returnType {
	case commandObject:
		if returnAsJSON {
			inst.cmd.SerializeObject = serializeCurrentFlowArray(object)
		} else {
			inst.cmd.Object = object
		}
		return inst.cmd
	case commandString:
		inst.cmd.String = inst.handleGetCommandString(command, object, parsedArgs)
		return inst.cmd
	default:

	}
	inst.cmd.Error = fmt.Errorf("unsupported return type: %s", returnType)
	return inst.cmd
}

type parsedCommand struct {
	ID       string `json:"id,omitempty"`
	Field    string `json:"field,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Write    string `json:"write,omitempty"`
	Value    string `json:"value,omitempty"`
	ReturnAs string `json:"returnAs"`
	IsQuery  bool   `json:"isQuery"`
}

const (
	commandJSON    = "json"
	commandPort    = "port"
	commandPorts   = "ports"
	commandNumber  = "number"
	commandString  = "string"
	commandObject  = "object"
	commandObjects = "objects"
	commandOutput  = "output"
	commandOutputs = "outputs"
	commandInput   = "input"
	commandInputs  = "inputs"
)

func commandReturnType(cmd *Command) (string, *parsedCommand) {
	var isQuery bool
	if cmd.Query != "" {
		isQuery = true
	}
	var isSet bool
	if cmd.CommandType == "set" {
		isSet = true
	}

	args := &parsedCommand{
		IsQuery: isQuery,
	}
	if v, ok := cmd.Args["name"]; ok {
		args.Name = v
	}
	if v, ok := cmd.Args["uuid"]; ok {
		args.UUID = v
	}
	if v, ok := cmd.Args["id"]; ok {
		args.ID = v
	}
	if v, ok := cmd.Args["field"]; ok {
		args.Field = v
	}
	if v, ok := cmd.Args["write"]; ok { // write in most cases will normally return nothing or an error
		args.Write = v
	}
	if v, ok := cmd.Args["value"]; ok {
		args.Value = v
	}
	if v, ok := cmd.Args["as"]; ok {
		args.ReturnAs = v
	}

	if cmd.Thing == commandObjects {
		if args.Write != "" {
			return commandString, args
		}
		return commandObjects, args
	}

	if cmd.Thing == commandObject {
		if args.Write != "" {
			return commandObject, args
		}
		return commandObject, args
	}

	if cmd.Thing == commandInputs || cmd.Thing == commandOutputs {
		if isSet { // update name, value of a port
			return commandString, args
		}
		return commandPorts, args
	}
	if cmd.Thing == commandInput || cmd.Thing == commandOutput {
		if isSet { // update name, value of a port
			return commandString, args
		}
		return commandPort, args
	}

	return "", args
}

// --------------- COMMANDS ----------------

func (inst *RuntimeImpl) handleGetManyAsObjects(cmd *Command, objects []Object, parsed *parsedCommand) []Object {
	var results []Object
	for _, object := range objects {
		result := inst.handleGetCommandObject(cmd, object, parsed)
		results = append(results, result)
	}
	return results
}

func (inst *RuntimeImpl) handleGetManyAsString(cmd *Command, objects []Object, parsed *parsedCommand) []string {
	var results []string
	for _, object := range objects {
		result := inst.handleGetCommandString(cmd, object, parsed)
		results = append(results, result)
	}
	return results
}

func (inst *RuntimeImpl) handleGetCommandObject(cmd *Command, object Object, parsed *parsedCommand) Object {
	switch strings.ToLower(cmd.Thing) {
	case "object":
		return object
	}
	return nil

}

func (inst *RuntimeImpl) handleGetCommandString(cmd *Command, object Object, parsed *parsedCommand) string {
	switch strings.ToLower(cmd.Thing) {
	case "uuid":
		return object.GetUUID()
	case "name":
		return object.GetName()
	}
	return ""

}

func (inst *RuntimeImpl) handleGetCommandForMultipleObjects(cmd *Command, objects []Object) any {
	var results []any
	for _, object := range objects {
		result := inst.handleGetCommand(cmd, object)
		results = append(results, result)
	}
	return results
}

func (inst *RuntimeImpl) handleSetCommandForMultipleObjects(cmd *Command, objects []Object) any {
	var results []any
	for _, object := range objects {
		result := inst.handleSetCommand(cmd, object)
		results = append(results, result)
	}
	return results
}

func (inst *RuntimeImpl) handleGetCommand(cmd *Command, object Object) any {
	switch strings.ToLower(cmd.Thing) {
	case "object":
		return object
	case "uuid":
		return object.GetUUID()
	case "name":
		return object.GetName()
	case "inputs":
		return getPortValues(object.GetInputs())
	case "input":
		return inst.handlePort(cmd, object, true)
	case "outputs":
		return getPortValues(object.GetOutputs())
	case "output":
		return inst.handlePort(cmd, object, false)
	default:
		return fmt.Errorf("unknown get command: %s", cmd.Thing)
	}

}

func (inst *RuntimeImpl) handlePort(cmd *Command, object Object, isInput bool) any {
	getID := cmd.GetArgsByKey("id")
	if getID == "" {
		return fmt.Errorf("failed to get value required from args :%s", "id")
	}
	var port *Port
	if isInput {
		port = object.GetInput(getID)
	} else {
		port = object.GetOutput(getID)
	}
	if port == nil {
		return fmt.Errorf("failed to find port by id: %s", getID)
	}
	get := cmd.GetArgsByKey("name")
	if get != "" {
		return port.GetValueDisplay()
	}
	get = cmd.GetArgsByKey("value")
	if get != "" {
		return port.GetValueDisplay()
	}
	get = cmd.GetArgsByKey("data")
	if get != "" {
		return port.GetValueDisplay()
	}
	port.DataDisplay = port.GetValueDisplay()
	return port
}

func (inst *RuntimeImpl) handleSetCommand(cmd *Command, object Object) any {
	pprint.PrintJSON(cmd)
	switch strings.ToLower(cmd.Thing) {
	case "object":
		field := cmd.GetArgsByKey("field")
		value := cmd.GetArgsByKey("value")
		if field == "" {
			return fmt.Errorf("failed to get value required from args :%s", "name")
		}
		if value == "" {
			return fmt.Errorf("failed to get value required from args :%s", "value")
		}
		if field == "name" {
			return object.SetName(value)
		}

	case "input":
		getID := cmd.GetArgsByKey("id")
		if getID == "" {
			return fmt.Errorf("failed to get value required from args :%s", "id")
		}
		var port *Port
		port = object.GetInput(getID)
		write := cmd.GetArgsByKey("write")
		if write != "" {
			if port.GetDataType() == priority.TypeFloat {
				f := convert.AnyToFloatPointer(write)
				if f == nil {
					return "was unable to parse value as type float"
				}
				err := port.Write(f)
				if err != nil {
					return err.Error()
				}
				return fmt.Sprintf("object: %s updated ok port: %s value: %s", object.GetName(), port.GetID(), write)
			}
		}

	default:
		return fmt.Errorf("unknown set command: %s", cmd.Thing)
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
	query, limit := extractAndRemoveLimit(query) // get the query limit
	orSegments := strings.Split(query, "OR")

	var allResults []Object
	for _, orSegment := range orSegments {
		orSegment = strings.TrimSpace(strings.Trim(orSegment, "()"))
		andConditions := strings.Split(orSegment, "AND") // Split the OR segment into AND conditions

		if len(andConditions) == 0 {
			continue
		}

		segmentResults := inst.Get() // Start with all objects for the first condition
		segmentResults = filterObjectsByCondition(segmentResults, andConditions)
		allResults = addUniqueMatches(allResults, segmentResults)
	}

	if limit >= 0 && len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults
}

func addUniqueMatches(allResults, segmentResults []Object) []Object {
	for _, match := range segmentResults {
		if !containsObject(allResults, match) {
			allResults = append(allResults, match)
		}
	}
	return allResults
}

func filterObjectsByCondition(segmentResults []Object, andConditions []string) []Object {
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

	return segmentResults
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

// ObjectConfig represents configuration for a object.
type ObjectConfig struct {
	ID                 string        `json:"id"`
	Info               *Info         `json:"info"`
	Inputs             []*Port       `json:"inputs"`
	Outputs            []*Port       `json:"outputs,omitempty"`
	Values             []*Port       `json:"values,omitempty"`
	Connections        []*Connection `json:"connections,omitempty"`
	Settings           *Settings     `json:"settings,omitempty"`
	Meta               *Meta         `json:"meta,omitempty"`
	Stats              *ObjectStats  `json:"stats,omitempty"`
	WasUpdated         bool          `json:"wasUpdated,omitempty"`
	dontRecreateObject bool
}

func SerializeCurrentFlowArray(objects []Object) []*ObjectConfig {
	var serializedObjects []*ObjectConfig
	for _, object := range objects {
		serializedObjects = append(serializedObjects, serializeCurrentFlowArray(object))
	}
	return serializedObjects
}

func serializeCurrentFlowArray(object Object) *ObjectConfig {

	meta := object.GetMeta()
	if meta == nil {
		meta = &Meta{
			Position: Position{
				PositionY: 0,
				PositionX: 0,
			},
		}
	}
	objectConfig := &ObjectConfig{
		ID:          object.GetID(),
		Info:        object.GetInfo(),
		Inputs:      getPortValues(object.GetInputs()),
		Outputs:     getPortValues(object.GetOutputs()),
		Connections: object.GetConnections(),
		Settings:    object.GetSettings(),
		Stats:       object.GetStats(),
		Meta:        meta,
	}
	return objectConfig
}

func getPortValues(ports []*Port) []*Port {
	for _, port := range ports {
		if port.GetValue() != nil {
			port.DataDisplay = port.GetValueDisplay()
		}
	}
	return ports
}
