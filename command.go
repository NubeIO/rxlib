package rxlib

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"regexp"
	"strconv"
	"strings"
)

type ExtendedCommand struct {
	*runtime.Command
}

func NewCommand() *ExtendedCommand {
	return &ExtendedCommand{
		Command: &runtime.Command{
			Data: make(map[string]string),
		},
	}
}

func CommandPing() *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "ping", "", "", false)
	c.Key = "ping"
	return c
}

func GetSerializeObjectByUUID(value string) *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "object", "uuid", value, true)
	return c
}

func GetObjectByUUID(value string) *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "object", "uuid", value, false)
	return c
}

func GetSerializeObjectByName(value string) *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "object", "name", value, true)
	return c
}

func GetObjectByName(value string) *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "object", "name", value, false)
	return c
}

func GetObject111111(value string) *ExtendedCommand {
	c := NewCommand()
	c.buildCommand("get", "objects", "name", value, false)
	return c
}

// QueryObjectsByField
// eg; QueryObjectByField("category", "math", 1, false)
func (c *ExtendedCommand) QueryObjectsByField(field, value string, limit int, asJSON bool) *ExtendedCommand {
	if limit > 0 {
		c.Query = fmt.Sprintf("objects:%s == %s limit:%d", field, value, limit)
	} else {
		c.Query = fmt.Sprintf("objects:%s == %s", field, value)
	}
	c.buildCommand("get", "objects", "type", value, asJSON)
	return c
}

func (c *ExtendedCommand) GetArgsByIndex(i int) string {
	if len(c.Args) > i {
		return c.Args[i]
	}
	return ""
}

// GetArgsByKey retrieves the value of a specific key in the Args map.
func (c *ExtendedCommand) GetArgsByKey(key string) string {
	return c.Data[key]
}

// GetArgsKeys retrieves all keys and values from the Args map.
func (c *ExtendedCommand) GetArgsKeys() (keys, values []string) {
	for k, v := range c.Data {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func (c *ExtendedCommand) GetName() string {
	return c.GetArgsByKey("name")
}

func (c *ExtendedCommand) GetUUID() string {
	return c.GetArgsByKey("uuid")
}

func (c *ExtendedCommand) GetID() string {
	return c.GetArgsByKey("id")
}

func (c *ExtendedCommand) GetField() string {
	return c.GetArgsByKey("field")
}

func (c *ExtendedCommand) GetKey() string {
	if c.Key == "" {
		return c.GetArgsByKey("key")
	}
	return c.Key
}

func (c *ExtendedCommand) GetObjectByUUID(value string, asJSON bool) *ExtendedCommand {
	c.buildCommand("get", "object", "uuid", value, asJSON)
	return c
}

func (c *ExtendedCommand) buildCommand(commandType, thing, fieldName, fieldValue string, asJSON bool) *ExtendedCommand {
	c.Args = append(c.Args, commandType)
	c.Args = append(c.Args, thing)
	if fieldName != "" {
		c.Data[fieldName] = fieldValue
	}
	if asJSON {
		c.Data["as"] = "json"
	}
	return c
}
func (c *ExtendedCommand) Parse(cmdSting string) (*ExtendedCommand, error) {
	argMap := make(map[string]string)
	var posArgs []string // Initialize as an empty slice

	// Split the string by spaces, but keep quoted strings together
	parts := splitArgs(cmdSting)
	for _, part := range parts {
		if strings.HasPrefix(part, "--") {
			// Parse key-value pair
			part = part[2:] // Remove the "--" prefix
			var keyValue []string
			if strings.Contains(part, "=") {
				keyValue = strings.SplitN(part, "=", 2)
			} else if strings.Contains(part, ":") {
				keyValue = strings.SplitN(part, ":", 2)
			}

			if len(keyValue) == 2 {
				key := keyValue[0]
				value := strings.Trim(keyValue[1], "\"")
				argMap[key] = value
			}
		} else {
			// Parse positional argument
			trimmedPart := strings.Trim(part, "\"")
			posArgs = append(posArgs, trimmedPart)
		}
	}
	if len(posArgs) == 1 {
		parts := splitCamelCase(posArgs[0])
		if len(parts) > 1 {
			cmdType := strings.ToLower(parts[0])
			things := strings.ToLower(parts[1])
			posArgs = nil
			posArgs = append(posArgs, cmdType)
			posArgs = append(posArgs, things)
		}
	}

	out := &ExtendedCommand{
		Command: &runtime.Command{
			Args: posArgs,
			Data: argMap,
		},
	}
	out.Key = argMap["key"]
	return out, nil
}

func UnmarshalCommand(payload any) (*ExtendedCommand, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var data *ExtendedCommand
	err = json.Unmarshal(marshal, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalCommandResponse(payload any) (*runtime.CommandResponse, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var data *runtime.CommandResponse
	err = json.Unmarshal(marshal, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func splitArgs(input string) []string {
	var parts []string
	var currentPart []rune
	inQuotes := false

	for _, r := range input {
		switch r {
		case ' ':
			if inQuotes {
				currentPart = append(currentPart, r)
			} else {
				if len(currentPart) > 0 {
					parts = append(parts, string(currentPart))
					currentPart = []rune{}
				}
			}
		case '"':
			inQuotes = !inQuotes
			currentPart = append(currentPart, r)
		default:
			currentPart = append(currentPart, r)
		}
	}

	if len(currentPart) > 0 {
		parts = append(parts, string(currentPart))
	}

	return parts
}

type CommandType string

const (
	CommandTypeGet    CommandType = "get"
	CommandTypeSet    CommandType = "set"
	CommandTypeDelete CommandType = "delete"
)

func (p *ParsedCommand) GetCommandType() string {
	return p.CommandType
}

func (p *ParsedCommand) IsSet() bool {
	if p.GetCommandType() == "set" {
		return true
	}
	return false
}

func (p *ParsedCommand) IsGet() bool {
	if p.GetCommandType() == "get" {
		return true
	}
	return false
}

func (p *ParsedCommand) IsRun() bool {
	if p.GetCommandType() == "run" {
		return true
	}
	return false
}

func (p *ParsedCommand) NameUUID() bool {
	if p.GetID() != "" {
		return true
	}
	if p.GetName() != "" {
		return true
	}
	if p.GetUUID() != "" {
		return true
	}
	return false
}

func (p *ParsedCommand) GetField() string {
	return p.Field
}

func (p *ParsedCommand) GetKey() string {
	return p.Key
}

func (p *ParsedCommand) IsFieldPort() bool {
	if p.Field == "input" {
		return true
	}
	if p.Field == "inputs" {
		return true
	}
	if p.Field == "output" {
		return true
	}
	if p.Field == "outputs" {
		return true
	}
	return false
}

func (p *ParsedCommand) GetThing() string {
	return p.Thing
}

func (p *ParsedCommand) ThingIsPorts() bool {
	if p.Thing == "input" {
		return true
	}
	if p.Thing == "inputs" {
		return true
	}
	if p.Thing == "output" {
		return true
	}
	if p.Thing == "outputs" {
		return true
	}
	return false
}

func (p *ParsedCommand) ThingIsObject() bool {
	if p.Thing == "object" {
		return true
	}
	if p.Thing == "objects" {
		return true
	}
	return false
}

func (p *ParsedCommand) GetQuery() string {
	return p.Query
}

func (p *ParsedCommand) GetCategory() string {
	return p.Category
}

func (p *ParsedCommand) GetID() string {
	return p.ID
}

func (p *ParsedCommand) GetReturnAs() string {
	return p.ReturnAs
}

func (p *ParsedCommand) GetChilds() bool {
	return p.Childs
}

func (p *ParsedCommand) GetPagination() bool {
	return p.Pagination
}

func (p *ParsedCommand) GetPaginationPageSize() int {
	return p.PageSize
}

func (p *ParsedCommand) GetPaginationPageNumber() int {
	return p.PageNumber
}

func (p *ParsedCommand) GetName() string {
	return p.Name
}

func (p *ParsedCommand) GetUUID() string {
	return p.UUID
}

func (p *ParsedCommand) SetReturnAsIfNil(s string) {
	p.ReturnAs = s
}

func splitCamelCase(s string) []string {
	re := regexp.MustCompile(`[a-z]+|[A-Z][a-z]*|[A-Z]+`)
	words := re.FindAllString(s, -1)
	return words
}

func (c *ExtendedCommand) ParseCommandsArgs(cmd *ExtendedCommand) (*ParsedCommand, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command can not be empty")
	}
	args := &ParsedCommand{}

	args.CommandType = cmd.GetArgsByIndex(0)
	args.Thing = cmd.GetArgsByIndex(1)

	if v, ok := cmd.Data["name"]; ok {
		args.Name = v
	}
	if v, ok := cmd.Data["uuid"]; ok {
		args.UUID = v
	}
	if v, ok := cmd.Data["id"]; ok {
		args.ID = v
	}
	if v, ok := cmd.Data["category"]; ok {
		args.Category = v
	}
	if v, ok := cmd.Data["query"]; ok {
		args.Query = v
	}
	if v, ok := cmd.Data["field"]; ok {
		args.Field = v
	}
	if v, ok := cmd.Data["write"]; ok { // write in most cases will normally return nothing or an error
		args.Write = v
	}
	if v, ok := cmd.Data["value"]; ok {
		args.Value = v
	}
	if v, ok := cmd.Data["as"]; ok {
		args.ReturnAs = v
	}
	if v, ok := cmd.Data["childs"]; ok {
		args.Childs = stringToBool(v)
	}
	// ObjectPagination
	// pageNumber, pageSize
	if v, ok := cmd.Data["pagination"]; ok {
		args.Pagination = stringToBool(v)
	}
	if v, ok := cmd.Data["pageNumber"]; ok {
		if args.Pagination {
			args.PageNumber = stringToInt(v)
		}
	}
	if v, ok := cmd.Data["pageSize"]; ok {
		if args.Pagination {
			args.PageSize = stringToInt(v)
		}
	}
	switch args.GetThing() {
	case "ping":
		return args, nil
	case commandCommand:
		return args, nil
	case commandObjects, commandObject:
		return args, nil
	case commandInputs, commandOutputs, commandInput, commandOutput:
		if args.IsSet() || (args.IsGet() && args.GetField() == "data") || (args.IsGet() && args.GetField() != "") {
			//args.ReturnAs = commandString
			args.SetReturnAsIfNil(commandString)
			return args, nil
		}
		return args, nil
	default:
		return args, fmt.Errorf("failed to find a valid type, get input, getInput or setInput")
	}
}

func IsObject(parsed *ParsedCommand) bool {
	if parsed.Thing == "object" {
		return true
	}
	return false
}

func IsObjects(parsed *ParsedCommand) bool {
	if parsed.Thing == "objects" {
		return true
	}
	return false
}

func IsInput(parsed *ParsedCommand) bool {
	if parsed.Thing == "inputs" || parsed.Thing == "input" {
		return true
	}
	return false
}

func IsPort(parsed *ParsedCommand) bool {
	if parsed.Thing == "inputs" || parsed.Thing == "input" {
		return true
	}
	if parsed.Thing == "outputs" || parsed.Thing == "output" {
		return true
	}
	return false
}

func ConvertCommand(c *runtime.Command) *ExtendedCommand {
	out := &ExtendedCommand{
		Command: &runtime.Command{
			TargetGlobalID:   c.GetTargetGlobalID(),
			SenderGlobalID:   c.GetSenderGlobalID(),
			SenderObjectUUID: c.GetSenderObjectUUID(),
			TransactionUUID:  c.GetTransactionUUID(),
			Key:              c.GetKey(),
			Args:             c.GetArgs(),
			Data:             c.GetData(),
		},
	}
	return out
}

func stringToBool(input interface{}) bool {
	switch v := input.(type) {
	case string:
		lowerInput := strings.ToLower(v)
		return lowerInput == "true"
	case bool:
		return v
	default:
		return false
	}
}

func stringToInt(input interface{}) int {
	switch v := input.(type) {
	case string:
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return num
	case int:
		return v
	default:
		return 0
	}
}
