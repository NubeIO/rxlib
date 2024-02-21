package rxlib

import (
	"regexp"
	"strings"
)

func splitCamelCase(s string) []string {
	re := regexp.MustCompile(`[a-z]+|[A-Z][a-z]*|[A-Z]+`)
	words := re.FindAllString(s, -1)
	return words
}

type Command struct {
	SenderID string            `json:"senderID"` // if sent from another ROS instance
	Key      string            `json:"key,omitempty"`
	Query    string            `json:"query,omitempty"`
	Args     []string          `json:"args,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
	Body     any               `json:"body,omitempty"`
}

func NewCommand() *Command {
	return &Command{
		Data: make(map[string]string),
	}
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
