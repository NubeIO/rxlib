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
		return nil, fmt.Errorf("failed to find vaild command convertion, try setName, getName, getInput")
	}
	commandType := strings.ToLower(commandParts[0])
	commandName := strings.ToLower(commandParts[1])
	if commandType == "set" || commandType == "write" {
		commandType = "set"
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
	Name        string   `json:"name,omitempty"`
	CommandName string   `json:"commandName,omitempty"`
	Query       string   `json:"query,omitempty"`
	Field       string   `json:"field,omitempty"`
	FieldEntry  string   `json:"fieldEntry,omitempty"`
	Args        []string `json:"args,omitempty"`
}

type PredefinedCommandBuilder struct {
	command *PredefinedCommand
}

func NewPredefinedCommandBuilder() *PredefinedCommandBuilder {
	return &PredefinedCommandBuilder{command: &PredefinedCommand{}}
}

func (b *PredefinedCommandBuilder) SetName(name string) *PredefinedCommandBuilder {
	b.command.Name = name
	return b
}

func (b *PredefinedCommandBuilder) SetCommandName(commandName string) *PredefinedCommandBuilder {
	b.command.CommandName = commandName
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
		case "-name":
			b.command.Field = "name"
			b.command.FieldEntry = value
		case "-id":
			b.command.Args = append(b.command.Args, "id", value)
		case "-write":
			b.command.Args = append(b.command.Args, "write", value)
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
		CommandName: pc.CommandName,
		Query:       pc.Query,
		Field:       pc.Field,
		FieldEntry:  pc.FieldEntry,
		Args:        pc.Args,
	}
	return cmd
}

type CommandBuilder struct {
	command *Command
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		command: &Command{},
	}
}

func (b *CommandBuilder) SetType(commandType CommandType) *CommandBuilder {
	b.command.CommandType = commandType
	return b
}

func (b *CommandBuilder) SetCommandName(name string) *CommandBuilder {
	b.command.CommandName = name
	return b
}

func (b *CommandBuilder) SetQuery(query string) *CommandBuilder {
	b.command.Query = query
	return b
}

func (b *CommandBuilder) SetField(field, entry string) *CommandBuilder {
	b.command.Field = field
	b.command.FieldEntry = entry
	return b
}

func (b *CommandBuilder) AddArg(arg string) *CommandBuilder {
	b.command.Args = append(b.command.Args, arg)
	return b
}

func (b *CommandBuilder) Build() *Command {
	return b.command
}

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

type Command struct {
	CommandType CommandType `json:"type"`            // get
	CommandName string      `json:"name"`            // -cmd:inputValues
	Query       string      `json:"query,omitempty"` //  -query:(objects:name == math-add-2) OR (objects:name == math-add-1)
	Field       string      `json:"field,omitempty"` // -uuid
	FieldEntry  string      `json:"entry,omitempty"` // abc
	Args        []string    `json:"args,omitempty"`  // all | in1
}

func NewCommand(c *Command) *Command {
	return c
}

func (c *Command) ToString() string {
	var fieldPart, queryPart string

	if c.Field != "" && c.FieldEntry != "" {
		fieldPart = fmt.Sprintf("-%s:%s", c.Field, c.FieldEntry)
	}

	if c.Query != "" {
		queryPart = fmt.Sprintf("-query:%s", c.Query)
	}

	return fmt.Sprintf("%s -cmd:%s %s %s -args: %s",
		c.CommandType, c.CommandName, queryPart, fieldPart, strings.Join(c.Args, ", "))
}

type CommandType string

const (
	CommandTypeGet    CommandType = "get"
	CommandTypeSet    CommandType = "set"
	CommandTypeDelete CommandType = "delete"
)
