package rxlib

import (
	"fmt"
	"regexp"
	"strings"
)

type CommandParse struct{}

func NewCommandParse() *CommandParse {
	return &CommandParse{}
}

func (c *CommandParse) Parse(command string) (*Command, error) {
	parts := strings.Fields(command)
	if len(parts) < 1 {
		return nil, fmt.Errorf("incorrect command")
	}
	cmd := c.Builder(command)
	return cmd, nil
}

func (c *CommandParse) Builder(command string) *Command {
	return NewPredefinedCommandBuilder().
		AddArgs(command).
		Build().
		ToCommand()
}

func splitCamelCase(s string) []string {
	re := regexp.MustCompile(`[a-z]+|[A-Z][a-z]*|[A-Z]+`)
	words := re.FindAllString(s, -1)
	return words
}

type PredefinedCommand struct {
	Query string            `json:"query,omitempty"`
	Args  map[string]string `json:"args,omitempty"`
}

type PredefinedCommandBuilder struct {
	command *PredefinedCommand
}

func NewPredefinedCommandBuilder() *PredefinedCommandBuilder {
	return &PredefinedCommandBuilder{command: &PredefinedCommand{
		Args: make(map[string]string),
	}}
}

func (b *PredefinedCommandBuilder) AddArgs(args string) *PredefinedCommandBuilder {
	var parts []string
	words := strings.Fields(args)
	if len(words) < 2 {
		return nil
	}

	commandType := words[0]
	var commandTypeSet bool
	if commandType == "set" {
		commandTypeSet = true
	}
	var commandTypeGet bool
	if commandType == "get" {
		commandTypeGet = true
	}

	if !commandTypeSet {
		if !commandTypeGet {
			words = splitCamelCase(args)
			commandType = words[0]
		}
	}

	commandType = strings.ToLower(commandType) // test for case getInputs, GetInputs
	thing := strings.ToLower(words[1])
	if commandType != "set" {
		if commandType != "get" {
			return nil
		}
	}

	if strings.Contains(args, "--query:") {
		// Split the string at the query
		querySplit := strings.SplitN(args, "--query:", 2)
		// Add the parts before the query
		parts = append(parts, strings.Fields(querySplit[0])...)
		// Handle the query separately
		if len(querySplit) > 1 {
			queryEndIndex := strings.Index(querySplit[1], " -")
			queryValue := querySplit[1]
			if queryEndIndex != -1 {
				// Query is followed by other arguments
				queryValue = querySplit[1][:queryEndIndex]
				// Add the parts after the query
				remainingArgs := querySplit[1][queryEndIndex+1:]
				parts = append(parts, parseArgsWithQuotes(remainingArgs)...)
			}
			// Ensure the query is enclosed in quotes
			if !strings.HasPrefix(queryValue, "\"") {
				queryValue = "\"" + queryValue + "\""
			}
			parts = append(parts, "-query:"+queryValue)
		}
	} else {
		parts = parseArgsWithQuotes(args)
	}

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		if commandType == "set" {
			b.command.Args["cmd"] = "set"
		} else {
			b.command.Args["cmd"] = "get"
		}
		b.command.Args["thing"] = thing
		key, value := kv[0], strings.Trim(kv[1], "\"")
		switch key {
		case "-query":
			b.command.Query = value
		case "-name":
			b.command.Args["name"] = value
		case "-uuid":
			b.command.Args["uuid"] = value
		case "-id":
			b.command.Args["id"] = value
		case "-write":
			b.command.Args["write"] = value
		case "-value":
			b.command.Args["value"] = value
		case "-field":
			b.command.Args["field"] = value
		case "-as":
			b.command.Args["as"] = value
		}
	}
	return b
}

func parseArgsWithQuotes(args string) []string {
	var parts []string
	var inQuotes bool
	var currentArg string

	for _, r := range args {
		switch r {
		case '"':
			inQuotes = !inQuotes
			currentArg += string(r)
		case ' ':
			if inQuotes {
				currentArg += string(r)
			} else {
				if currentArg != "" {
					parts = append(parts, currentArg)
					currentArg = ""
				}
			}
		case '-':
			if !inQuotes && currentArg != "" {
				parts = append(parts, currentArg)
				currentArg = ""
			}
			currentArg += string(r)
		default:
			currentArg += string(r)
		}
	}
	if currentArg != "" {
		parts = append(parts, currentArg)
	}

	return parts
}

func (b *PredefinedCommandBuilder) Build() *PredefinedCommand {
	if b == nil {
		return nil
	}
	if b.command == nil {
		return nil
	}
	return b.command
}

func (pc *PredefinedCommand) ToCommand() *Command {
	if pc == nil {
		return nil
	}
	cmd := &Command{
		Query: pc.Query,
	}
	cmd.Args = make(map[string]string)
	cmd.Args = pc.Args
	return cmd
}

type Command struct {
	Query string            `json:"query,omitempty"` //  -query:(objects:name == math-add-2) OR (objects:name == math-add-1)
	Args  map[string]string `json:"args,omitempty"`
}
type CommandOld2 struct {
	CommandType CommandType       `json:"type"`            // get
	Thing       string            `json:"thing"`           // -cmd:inputValues
	Query       string            `json:"query,omitempty"` //  -query:(objects:name == math-add-2) OR (objects:name == math-add-1)
	Field       string            `json:"field,omitempty"` // -uuid
	FieldEntry  string            `json:"entry,omitempty"` // abc
	Args        map[string]string `json:"args,omitempty"`
}

func NewCommand() *Command {
	return &Command{
		Args: make(map[string]string),
	}
}

// GetArgsByKey retrieves the value of a specific key in the Args map.
func (c *Command) GetArgsByKey(key string) string {
	return c.Args[key]
}

// GetArgsKeys retrieves all keys and values from the Args map.
func (c *Command) GetArgsKeys() (keys, values []string) {
	for k, v := range c.Args {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func (c *Command) ToString() string {
	var fieldPart, queryPart string

	if c.Query != "" {
		queryPart = fmt.Sprintf("-query:%s", c.Query)
	}

	var args []string
	for key, value := range c.Args {
		args = append(args, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	return fmt.Sprintf("-cmd:%s %s %s -args: %s",
		queryPart, fieldPart, strings.Join(args, ", "))
}

type CommandType string

const (
	CommandTypeGet    CommandType = "get"
	CommandTypeSet    CommandType = "set"
	CommandTypeDelete CommandType = "delete"
)

type CommandOld struct {
	NameSpace  string `json:"nameSpace,omitempty"`  // name space is like this action.<plugin>.<objectUUID>.<commandName>
	ObjectUUID string `json:"objectUUID,omitempty"` // or use UUID
	Key        string `json:"key,omitempty"`
	Body       any    `json:"body,omitempty"`
}

type Action struct {
	Plugin      string
	ObjectName  string
	CommandName string
	Body        any
}
