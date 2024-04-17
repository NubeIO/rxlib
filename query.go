package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/convert"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"strconv"
	"strings"
)

func convertCommand(resp *CommandResponse) *runtime.CommandResponse {
	return &runtime.CommandResponse{
		SenderID:         resp.SenderID,
		Count:            int32(resp.Count),
		MapStrings:       resp.MapStrings,
		Number:           resp.Float,
		Boolean:          resp.Bool,
		Error:            resp.Error,
		ReturnType:       resp.ReturnType,
		Any:              resp.Any,
		Response:         convertCommands(resp.CommandResponse),
		SerializeObjects: resp.SerializeObjects,
		ObjectPagination: resp.ObjectPagination,
	}
}

func convertCommands(resp []*CommandResponse) []*runtime.CommandResponse {
	var out []*runtime.CommandResponse
	for _, response := range resp {
		out = append(out, convertCommand(response))
	}
	return out
}

func (inst *RuntimeImpl) Command(cmd *ExtendedCommand) *runtime.CommandResponse {
	resp := inst.CommandObject(cmd)
	return convertCommand(resp)

}

func (inst *RuntimeImpl) CommandObject(command *ExtendedCommand) *CommandResponse {
	inst.response = &CommandResponse{
		MapStrings: make(map[string]string),
	}
	if command == nil {
		inst.response.Error = fmt.Sprintf("command cannot be nil")
		return inst.response
	}
	inst.command = command
	parsedArgs, err := inst.command.ParseCommandsArgs(command)
	if err != nil {
		inst.response.Error = fmt.Sprintf("%v", err)
		return inst.response
	}

	inst.response.SenderID = command.SenderGlobalID
	inst.response.ReturnType = parsedArgs.GetReturnAs()
	inst.parsedCommand = parsedArgs

	switch parsedArgs.Thing {
	case "objects", "object", "command":
		return inst.handleObjects(parsedArgs)
	default:
		inst.response.Error = fmt.Sprintf("unknown command type: %s", parsedArgs.Thing)
		return inst.response
	}
}

func (inst *RuntimeImpl) handleObjects(parsedArgs *ParsedCommand) *CommandResponse {
	var objects []Object
	if parsedArgs.GetPagination() {
		pagination, err := inst.handlePagination(parsedArgs)
		if err != nil {
			objectsLen := len(objects)
			fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
			inst.response.Error = fmt.Sprintf("failed to find any objects")
			return inst.response
		}
		objects = pagination.Objects
		pagination.Objects = nil
		inst.response.ObjectPagination = &runtime.ObjectPagination{
			Count:      int32(pagination.Count),
			PageNumber: int32(pagination.PageNumber),
			PageSize:   int32(pagination.PageSize),
			TotalPages: int32(pagination.TotalPages),
			TotalCount: int32(pagination.TotalCount),
		}
		inst.handleReturnType(parsedArgs, objects)
		return inst.response
	}
	if parsedArgs.GetTree() { // handel object tree
		inst.response.ObjectTree = inst.GetTreeMapRoot()
		return inst.response
	}

	objects = inst.getObjects(parsedArgs)
	objectsLen := len(objects)
	if objectsLen == 0 {
		fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
		inst.response.Error = fmt.Sprintf("failed to find any objects")
		return inst.response
	}
	inst.handleCommandTypeObjects(parsedArgs, objects)
	inst.handleCommandTypePorts(parsedArgs, objects)
	inst.handleReturnType(parsedArgs, objects)
	fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
	return inst.response
}

func (inst *RuntimeImpl) getObjects(parsedArgs *ParsedCommand) []Object {
	if parsedArgs.GetQuery() != "" {
		return inst.query(parsedArgs.GetQuery())
	}
	return inst.handleNoQuery(parsedArgs)

}

func (inst *RuntimeImpl) handleNoQuery(parsedArgs *ParsedCommand) []Object {
	switch parsedArgs.GetThing() {
	case "object": // get a single object
		object := inst.handleGetObject(parsedArgs)
		if object != nil {
			return []Object{object}
		}
		return nil
	default: // objects, command
		if !parsedArgs.NameUUID() { // get objects eg; getObjects --as:json
			return inst.Get()
		}
		return inst.handleGetObjects(parsedArgs)
	}
}

func (inst *RuntimeImpl) handleGetObjects(parsedArgs *ParsedCommand) []Object {
	var objects []Object
	for _, object := range inst.Get() {
		if parsedArgs.GetName() == object.GetName() {
			objects = append(objects, object)
		}
		if parsedArgs.GetUUID() == object.GetUUID() {
			objects = append(objects, object)
		}
		if parsedArgs.GetID() == object.GetID() {
			objects = append(objects, object)
		}
		if parsedArgs.GetCategory() == object.GetCategory() {
			objects = append(objects, object)
		}
	}
	return objects
}

func (inst *RuntimeImpl) handleGetObject(parsedArgs *ParsedCommand) Object {
	var obj Object
	switch {
	case parsedArgs.GetUUID() != "":
		obj = inst.GetByUUID(parsedArgs.GetUUID())
	case parsedArgs.GetName() != "":
		obj = inst.GetFirstByName(parsedArgs.GetName())
	case parsedArgs.GetID() != "":
		obj = inst.GetFirstByID(parsedArgs.GetID())
	case parsedArgs.GetCategory() != "":
		obj = inst.GetFirstByID(parsedArgs.GetCategory())
	}
	return obj
}

func (inst *RuntimeImpl) handleCommandTypeObjects(parsedArgs *ParsedCommand, objects []Object) {
	if parsedArgs.IsFieldPort() {
		return
	}
	switch parsedArgs.GetCommandType() {
	case "run":
		if parsedArgs.GetThing() == "command" {
			parsedArgs.SetReturnAsIfNil("command")
			for _, object := range objects {
				inst.response.CommandResponse = append(inst.response.CommandResponse, object.CommandObject(inst.command))
			}
		}
	case "set":

	case "get":
		inst.response.MapStrings = inst.handleGetManyAsString(objects, parsedArgs)

	case "invoke":

	}
}

func (inst *RuntimeImpl) handleReturnType(parsedArgs *ParsedCommand, objects []Object) {
	switch parsedArgs.GetReturnAs() {
	case commandCount:
		if parsedArgs.ThingIsObject() {
			inst.response.Count = len(objects)
			inst.response.Objects = nil
		}
		if parsedArgs.ThingIsPorts() {
			inst.response.Count = len(inst.response.MapPorts)
			inst.response.MapPorts = nil
			inst.response.Objects = nil
		}
	case commandJSON:
		inst.response.SerializeObjects = SerializeCurrentFlowArray(objects)
		inst.response.Count = len(inst.response.SerializeObjects)
	case commandCommand:
		inst.response.Count = len(inst.response.CommandResponse)
		inst.response.Objects = nil
	case commandString:
		inst.response.Objects = nil
	case commandPorts:
		inst.response.Objects = nil
		inst.response.Count = len(inst.response.MapPorts)
	default:
		inst.response.Objects = objects
		inst.response.Count = len(objects)
	}
}

func (inst *RuntimeImpl) handleCommandTypePorts(parsedArgs *ParsedCommand, objects []Object) {
	//if !parsedArgs.IsFieldPort() {
	//	return
	//}
	//
	//inst.response.ReturnType = "ports"
	//parsedArgs.SetReturnAsIfNil(inst.response.ReturnType)
	//inst.response.MapPorts = make(map[string][]*Port)
	//switch parsedArgs.GetField() {
	//case "input":
	//	for _, object := range objects {
	//		port := object.GetInput(parsedArgs.GetID())
	//		fmt.Println(port, "port")
	//		if port != nil {
	//			inst.response.MapPorts[object.GetUUID()] = []*Port{port}
	//		}
	//	}
	//case "inputs":
	//	for _, object := range objects {
	//		inst.response.MapPorts[object.GetUUID()] = object.GetOutputs()
	//	}
	//case "output":
	//	for _, object := range objects {
	//		port := object.GetOutput(parsedArgs.GetID())
	//		if port != nil {
	//			inst.response.MapPorts[object.GetUUID()] = []*Port{port}
	//		}
	//	}
	//case "outputs":
	//	for _, object := range objects {
	//		inst.response.MapPorts[object.GetUUID()] = object.GetOutputs()
	//	}
	//}
}

func errMessage(message, returnType string, parsed *ParsedCommand) error {
	return fmt.Errorf("error-message: %s, type: %s, thing: %s, type-return: %s\n", message, parsed.CommandType, parsed.Thing, returnType)
}

// --------------- COMMANDS ----------------

func (inst *RuntimeImpl) handleGetManyAsObjects(objects []Object, parsed *ParsedCommand) []Object {
	var results []Object
	for _, object := range objects {
		result := inst.handleGetCommandObject(object, parsed)
		results = append(results, result)
	}
	return results
}

func (inst *RuntimeImpl) handleGetCommandObject(object Object, parsed *ParsedCommand) Object {
	switch strings.ToLower(parsed.Thing) {
	case "object":
		return object
	}
	return nil

}

func (inst *RuntimeImpl) handleGetManyAsString(objects []Object, parsed *ParsedCommand) map[string]string {
	results := make(map[string]string)
	for _, object := range objects {
		result := inst.handleGetCommandString(object, parsed)
		results[object.GetUUID()] = result
	}
	return results
}

func (inst *RuntimeImpl) handleGetCommandString(object Object, parsed *ParsedCommand) string {
	if IsPort(parsed) { // handle ports
		res, err := inst.handlePortString(object, parsed)
		if err != nil {
			return err.Error()
		}
		return res
	}

	switch strings.ToLower(parsed.Field) {
	case "id":
		return object.GetID()
	case "uuid":
		return object.GetUUID()
	case "name":
		return object.GetName()
	}
	return ""
}

func (inst *RuntimeImpl) handlePortString(object Object, parsed *ParsedCommand) (string, error) {
	//port, err := portCommon(object, IsInput(parsed), parsed)
	//if err != nil {
	//	return "", err
	//}
	//switch strings.ToLower(parsed.Field) {
	//case "values":
	//	pri := port.GetValueDisplay()
	//	if pri == nil {
	//		return "no values found", nil
	//	}
	//	out := fmt.Sprintf("DataType: %s, HighestPriority: %v RawValue %v", pri.DataType, pri.HighestPriority, pri.RawValue)
	//	return out, nil
	//case "value":
	//	return fmt.Sprint(port.GetHighestPriority()), nil
	//case "datatype":
	//	return string(port.GetDataType()), nil
	//case "uuid":
	//	return port.GetUUID(), nil
	//case "name":
	//	return port.GetName(), nil
	//}
	return "", nil
}

func (inst *RuntimeImpl) handleSetPorts(object Object, parsed *ParsedCommand) (string, error) {
	switch strings.ToLower(parsed.Thing) {
	case "input":
		getID := parsed.ID
		if getID == "" {
			return "", fmt.Errorf("failed to get value required from args :%s", "id")
		}
		var port *Port
		port = object.GetInput(getID)
		write := parsed.Write
		if write != "" {
			if port.GetDataType() == priority.TypeFloat {
				f := convert.AnyToFloatPointer(write)
				if f == nil {
					return "", fmt.Errorf("was unable to parse value as type float")
				}
				//err := port.Write(f)
				//if err != nil {
				//	return "", err
				//}
				return "", fmt.Errorf("object: %s updated ok port: %s value: %s", object.GetName(), port.GetID(), write)
			}
		}

	default:
		return "", fmt.Errorf("unknown set command: %s", parsed.Thing)
	}
	return "", fmt.Errorf("unknown get command")

}

func (inst *RuntimeImpl) handlePorts(object Object, parsed *ParsedCommand) ([]*Port, error) {
	isInputs := IsInput(parsed)
	if parsed.ID != "" { // handle a one port
		var out []*Port
		p, err := portCommon(object, isInputs, parsed)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
		return out, nil
	}
	var ports []*Port
	if isInputs {
		ports = object.GetInputs()
	} else {
		ports = object.GetOutputs()
	}
	if ports == nil {
		//return nil, fmt.Errorf("failed to find port by id: %s", getID)
		return nil, nil
	}
	return ports, nil
}

type ParsedCommand struct {
	CommandType string `json:"commandType"`
	Thing       string `json:"thing"`
	ID          string `json:"id,omitempty"`
	Field       string `json:"field,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Category    string `json:"category,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Write       string `json:"write,omitempty"`
	Value       string `json:"value,omitempty"`
	ReturnAs    string `json:"GetReturnAs"`
	Query       string `json:"query"`
	Key         string `json:"key"`
	Childs      bool   `json:"childs,omitempty"`
	Tree        bool   `json:"tree,omitempty"`
	Pagination  bool   `json:"pagination,omitempty"`
	PageNumber  int    `json:"pageNumber,omitempty"`
	PageSize    int    `json:"pageSize,omitempty"`
}

const (
	commandCount = "count"
	commandJSON  = "json"
	//commandPort    = "port"
	commandPorts   = "ports"
	commandNumber  = "number"
	commandString  = "string"
	commandCommand = "command"
	commandObject  = "object"
	commandObjects = "objects"
	commandOutput  = "output"
	commandOutputs = "outputs"
	commandInput   = "input"
	commandInputs  = "inputs"
)

func portCommon(object Object, isInput bool, parsed *ParsedCommand) (*Port, error) {
	getID := parsed.ID
	if getID == "" {
		return nil, fmt.Errorf("failed to get value required from args :%s", "id")
	}
	var port *Port
	if isInput {
		port = object.GetInput(getID)
	} else {
		port = object.GetOutput(getID)
	}
	if port == nil {
		var inputCount = len(object.GetInputs())
		var outputCount = len(object.GetOutputs())
		return nil, fmt.Errorf("failed to find port by id: %s object inputs count: %d object outputs count: %d", getID, inputCount, outputCount)
	}
	return port, nil
}

func (inst *RuntimeImpl) handleSetCommand(cmd *ExtendedCommand, object Object, parsed *ParsedCommand) any {
	switch strings.ToLower(parsed.Thing) {
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
		//var port *Port
		//port = object.GetInput(getID)
		write := cmd.GetArgsByKey("write")
		if write != "" {
			//if port.GetDataType() == priority.TypeFloat {
			//	f := convert.AnyToFloatPointer(write)
			//	if f == nil {
			//		return "was unable to parse value as type float"
			//	}
			//	err := port.Write(f)
			//	if err != nil {
			//		return err.Error()
			//	}
			//	return fmt.Sprintf("object: %s updated ok port: %s value: %s", object.GetName(), port.GetID(), write)
			//}
		}

	default:
		return fmt.Errorf("unknown set command: %s", parsed.Thing)
	}
	return fmt.Errorf("unknown get command")

}

func (inst *RuntimeImpl) query(query string) []Object {
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
	//switch attribute {
	//case "value":
	//	if value == "null" { // get all ports that are null/nil
	//		isNull := port.GetValue().IsNull()
	//		if isNull {
	//			return true
	//		}
	//		return false
	//	}
	//	compareValue, err := strconv.ParseFloat(value, 64)
	//	if err != nil {
	//		return false // Invalid comparison value, return false
	//	}
	//	portValue := port.GetValue().GetFloat()
	//	return compareFloat64(portValue, operator, compareValue)
	//default:
	//	return false // Invalid attribute, return false
	//}
	return false
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
	case "type":
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
