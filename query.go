package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"strings"
)

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
	case "whois":
		obj := inst.whois(parsedArgs)
		if obj != nil {
			inst.response.SerializeObjects = inst.SerializeObjects(false, []Object{obj})
			inst.response.Count = len(inst.response.SerializeObjects)
		}
		return inst.response
	case "values":
		return inst.handleValues(parsedArgs)
	case "objects", "object", "command":
		return inst.handleObjects(parsedArgs)
	default:
		inst.response.Error = fmt.Sprintf("unknown command type: %s", parsedArgs.Thing)
		return inst.response
	}
}

func (inst *RuntimeImpl) handleValues(parsedArgs *ParsedCommand) *CommandResponse {
	if parsedArgs.GetPortValues() && !parsedArgs.GetPagination() {
		return inst.handlePortValues(parsedArgs)
	}
	return nil
}

func (inst *RuntimeImpl) handleObjects(parsedArgs *ParsedCommand) *CommandResponse {
	if parsedArgs.GetPagination() {
		return inst.handlePaginationObjects(parsedArgs)
	}
	if parsedArgs.GetTree() {
		return inst.handleTreeObjects(parsedArgs)
	}
	return inst.handleRegularObjects(parsedArgs)
}

func (inst *RuntimeImpl) handlePaginationObjects(parsedArgs *ParsedCommand) *CommandResponse {
	pagination, err := inst.handlePagination(parsedArgs)
	if err != nil {
		return inst.handlePaginationError(parsedArgs)
	}
	inst.response.ObjectPagination = &runtime.ObjectPagination{
		Count:      int32(pagination.Count),
		PageNumber: int32(pagination.PageNumber),
		PageSize:   int32(pagination.PageSize),
		TotalPages: int32(pagination.TotalPages),
		TotalCount: int32(pagination.TotalCount),
	}
	objects := pagination.Objects
	pagination.Objects = nil
	inst.handleReturnType(parsedArgs, objects)
	return inst.response
}

func (inst *RuntimeImpl) handlePortValues(parsedArgs *ParsedCommand) *CommandResponse {
	parentUUID := parsedArgs.GetUUID()
	inst.response.PortValues = inst.GetObjectsValues(parentUUID)
	return inst.response
}

func (inst *RuntimeImpl) handlePortValuesPaginate(parsedArgs *ParsedCommand) *CommandResponse {
	objectUUID := parsedArgs.GetUUID()
	pageSize := parsedArgs.GetPaginationPageSize()
	pageNumber := parsedArgs.GetPaginationPageNumber()
	pagination := inst.GetObjectsValuesPaginate(objectUUID, pageNumber, pageSize)
	inst.response.ObjectPagination = &runtime.ObjectPagination{
		Count:      int32(pagination.Count),
		PageNumber: int32(pagination.PageNumber),
		PageSize:   int32(pagination.PageSize),
		TotalPages: int32(pagination.TotalPages),
		TotalCount: int32(pagination.TotalCount),
		PortValues: pagination.PortValues,
	}
	return inst.response
}

func (inst *RuntimeImpl) handlePaginationError(parsedArgs *ParsedCommand) *CommandResponse {
	objectsLen := len(inst.objects)
	fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
	inst.response.Error = fmt.Sprintf("failed to find any objects")
	return inst.response
}

func (inst *RuntimeImpl) handleTreeObjects(parsedArgs *ParsedCommand) *CommandResponse {
	if parsedArgs.GetAncestor() {
		if parsedArgs.GetChilds() {
			inst.response.AncestorObjectTree = inst.GetTreeChilds(parsedArgs.GetUUID())
			return inst.response
		}
		inst.response.AncestorObjectTree = inst.GetAncestorTreeByUUID(parsedArgs.GetUUID())
		return inst.response
	}
	inst.response.ObjectTree = inst.GetTreeMapRoot()
	return inst.response
}

func (inst *RuntimeImpl) handleRegularObjects(parsedArgs *ParsedCommand) *CommandResponse {
	objects := inst.getObjects(parsedArgs)
	objectsLen := len(objects)
	if objectsLen == 0 {
		fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
		inst.response.Error = fmt.Sprintf("failed to find any objects")
		return inst.response
	}
	inst.handleCommandTypeObjects(parsedArgs, objects)
	inst.handleReturnType(parsedArgs, objects)
	fmt.Printf("type: %s, thing: %s, return type: %s, objects effected: %d \n", parsedArgs.GetCommandType(), parsedArgs.GetThing(), parsedArgs.GetReturnAs(), objectsLen)
	return inst.response
}

func (inst *RuntimeImpl) getObjects(parsedArgs *ParsedCommand) []Object {
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
		inst.response.SerializeObjects = inst.SerializeObjects(parsedArgs.GetPortValues(), objects)
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
	Ancestor    bool   `json:"ancestor,omitempty"`
	PortValues  bool   `json:"portValues,omitempty"`
	Pagination  bool   `json:"pagination,omitempty"`
	PageNumber  int    `json:"pageNumber,omitempty"`
	PageSize    int    `json:"pageSize,omitempty"`

	Start  int  `json:"start,omitempty"`
	Finish int  `json:"finish,omitempty"`
	Global bool `json:"global,omitempty"`
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

func convertCommand(resp *CommandResponse) *runtime.CommandResponse {
	return &runtime.CommandResponse{
		SenderID:           resp.SenderID,
		Count:              int32(resp.Count),
		MapStrings:         resp.MapStrings,
		TypeFloat:          resp.Float,
		TypeBool:           resp.Bool,
		TypeError:          resp.Error,
		ReturnType:         resp.ReturnType,
		TypeByte:           resp.Byte,
		Response:           convertCommands(resp.CommandResponse),
		SerializeObjects:   resp.SerializeObjects,
		ObjectPagination:   resp.ObjectPagination,
		ObjectTree:         resp.ObjectTree,
		AncestorObjectTree: resp.AncestorObjectTree,
		PortValues:         resp.PortValues,
	}
}
