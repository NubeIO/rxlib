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
	name := parts[0]
	commandParts := splitCamelCase(name)
	if len(commandParts) != 2 {
		return nil, fmt.Errorf("failed to find valid command conversion, try setName, getName, getInput")
	}
	commandType := strings.ToLower(commandParts[0])
	commandName := strings.ToLower(commandParts[1])
	if commandType == "set" || commandType == "write" {
		commandType = "set"
	} else if commandType == "get" {
		// ok
	} else {
		return nil, fmt.Errorf("incorrect command for set, try set or write eg; setName, writeName")
	}

	cmd := c.Builder(commandType, commandName, command)
	return cmd, nil
}

func (c *CommandParse) Builder(commandType, commandName, command string) *Command {
	return NewPredefinedCommandBuilder().
		SetName(commandType).
		SetCommandName(commandName).
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
	Name       string            `json:"name,omitempty"`
	Thing      string            `json:"thing,omitempty"`
	Query      string            `json:"query,omitempty"`
	Field      string            `json:"field,omitempty"`
	FieldEntry string            `json:"fieldEntry,omitempty"`
	Args       map[string]string `json:"args,omitempty"`
}

type PredefinedCommandBuilder struct {
	command *PredefinedCommand
}

func NewPredefinedCommandBuilder() *PredefinedCommandBuilder {
	return &PredefinedCommandBuilder{command: &PredefinedCommand{
		Args: make(map[string]string),
	}}
}

func (b *PredefinedCommandBuilder) SetName(name string) *PredefinedCommandBuilder {
	b.command.Name = name
	return b
}

func (b *PredefinedCommandBuilder) SetCommandName(commandName string) *PredefinedCommandBuilder {
	b.command.Thing = commandName
	return b
}

func (b *PredefinedCommandBuilder) AddArgs(args string) *PredefinedCommandBuilder {
	var parts []string
	if strings.Contains(args, "-query:") {
		// Split the string at the query
		querySplit := strings.SplitN(args, "-query:", 2)
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
			continue // Invalid argument format
		}

		key, value := kv[0], strings.Trim(kv[1], "\"")
		switch key {
		case "-query":
			b.command.Query = value
		case "-returnType":
			b.command.Field = "returnType"
			b.command.FieldEntry = value
		case "-name":
			b.command.Field = "name"
			b.command.FieldEntry = value
		case "-uuid":
			b.command.Field = "uuid"
			b.command.FieldEntry = value
		case "-id":
			b.command.Args["id"] = value
		case "-write":
			b.command.Args["write"] = value
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
	return b.command
}

func (pc *PredefinedCommand) ToCommand() *Command {
	cmd := &Command{
		CommandType: CommandType(pc.Name),
		Thing:       pc.Thing,
		Query:       pc.Query,
		Field:       pc.Field,
		FieldEntry:  pc.FieldEntry,
	}
	cmd.Args = make(map[string]string)
	cmd.Args = pc.Args
	return cmd
}

type Command struct {
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

	if c.Field != "" && c.FieldEntry != "" {
		fieldPart = fmt.Sprintf("-%s:%s", c.Field, c.FieldEntry)
	}

	if c.Query != "" {
		queryPart = fmt.Sprintf("-query:%s", c.Query)
	}

	var args []string
	for key, value := range c.Args {
		args = append(args, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	return fmt.Sprintf("%s -cmd:%s %s %s -args: %s",
		c.CommandType, c.Thing, queryPart, fieldPart, strings.Join(args, ", "))
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
