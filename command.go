package rxlib

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func splitCamelCase(s string) []string {
	re := regexp.MustCompile(`[a-z]+|[A-Z][a-z]*|[A-Z]+`)
	words := re.FindAllString(s, -1)
	return words
}

type Command struct {
	SenderGlobalID   string            `json:"senderGlobalID,omitempty"`   // if sent from another ROS instance
	SenderObjectUUID string            `json:"senderObjectUUID,omitempty"` // if sent from another ROS instance
	Key              string            `json:"key,omitempty"`
	Query            string            `json:"query,omitempty"`
	Args             []string          `json:"args,omitempty"`
	Data             map[string]string `json:"data,omitempty"`
	Body             any               `json:"body,omitempty"`
}

func NewCommand() *Command {
	return &Command{
		Data: make(map[string]string),
	}
}

func (c *Command) GetObjectByName(value string, asJSON bool) *Command {
	c.buildCommand("get", "object", "name", value, asJSON)
	return c
}

// QueryObjectsByField
// eg; QueryObjectByField("category", "math", 1, false)
func (c *Command) QueryObjectsByField(field, value string, limit int, asJSON bool) *Command {
	if limit > 0 {
		c.Query = fmt.Sprintf("objects:%s == %s limit:%d", field, value, limit)
	} else {
		c.Query = fmt.Sprintf("objects:%s == %s", field, value)
	}
	c.buildCommand("get", "objects", "type", value, asJSON)
	return c
}

func (c *Command) GetObjectByUUID(value string, asJSON bool) *Command {
	c.buildCommand("get", "object", "uuid", value, asJSON)
	return c
}

func (c *Command) buildCommand(commandType, thing, fieldName, fieldValue string, asJSON bool) *Command {
	c.Args = append(c.Args, commandType)
	c.Args = append(c.Args, thing)
	c.Data[fieldName] = fieldValue
	if asJSON {
		c.Data["as"] = "json"
	}
	return c
}
func (c *Command) Parse(cmdSting string) (*Command, error) {
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
	return &Command{Args: posArgs, Data: argMap}, nil
}

func UnmarshalCommand(payload any) (*Command, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var data *Command
	err = json.Unmarshal(marshal, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalCommandResponse(payload any) (*CommandResponse, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var data *CommandResponse
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

func (c *Command) GetArgsByIndex(i int) string {
	if len(c.Args) > i {
		return c.Args[i]
	}
	return ""
}

// GetArgsByKey retrieves the value of a specific key in the Args map.
func (c *Command) GetArgsByKey(key string) string {
	return c.Data[key]
}

// GetArgsKeys retrieves all keys and values from the Args map.
func (c *Command) GetArgsKeys() (keys, values []string) {
	for k, v := range c.Data {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

type CommandType string

const (
	CommandTypeGet    CommandType = "get"
	CommandTypeSet    CommandType = "set"
	CommandTypeDelete CommandType = "delete"
)

func (p *parsedCommand) getCommandType() string {
	return p.CommandType
}

func (p *parsedCommand) isSet() bool {
	if p.getCommandType() == "set" {
		return true
	}
	return false
}

func (p *parsedCommand) isGet() bool {
	if p.getCommandType() == "get" {
		return true
	}
	return false
}

func (p *parsedCommand) getField() string {
	return p.Field
}

func (p *parsedCommand) isFieldPort() bool {
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

func (p *parsedCommand) getThing() string {
	return p.Thing
}

func (p *parsedCommand) getQuery() string {
	return p.Query
}

func (p *parsedCommand) getID() string {
	return p.ID
}

func (p *parsedCommand) returnAs() string {
	return p.ReturnAs
}

func (p *parsedCommand) getName() string {
	return p.Name
}

func (p *parsedCommand) getUUID() string {
	return p.UUID
}

func commandReturnType(cmd *Command) (*parsedCommand, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command can not be empty")
	}
	args := &parsedCommand{}

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
	if v, ok := cmd.Data["type"]; ok {
		args.Type = v
	}

	switch args.getThing() {
	case commandObjects, commandObject:
		return args, nil
	case commandInputs, commandOutputs, commandInput, commandOutput:
		if args.isSet() || (args.isGet() && args.getField() == "data") || (args.isGet() && args.getField() != "") {
			args.ReturnAs = commandString
			return args, nil
		}
		return args, nil
	default:
		return args, fmt.Errorf("failed to find a valid type, get input, getInput or setInput")
	}
}

func isObject(parsed *parsedCommand) bool {
	if parsed.Thing == "object" {
		return true
	}
	return false
}

func isObjects(parsed *parsedCommand) bool {
	if parsed.Thing == "objects" {
		return true
	}
	return false
}

func isInput(parsed *parsedCommand) bool {
	if parsed.Thing == "inputs" || parsed.Thing == "input" {
		return true
	}
	return false
}

func isPort(parsed *parsedCommand) bool {
	if parsed.Thing == "inputs" || parsed.Thing == "input" {
		return true
	}
	if parsed.Thing == "outputs" || parsed.Thing == "output" {
		return true
	}
	return false
}
