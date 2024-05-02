package jsonutils

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/tidwall/gjson"
)

// JSON interface defines methods for working with JSON.
type JSON interface {
	// Compare compares two JSON strings based on the provided fields and options.
	Compare(j1, j2 string, fields []string, fieldsAsException []string, returnMatching, returnNonMatching bool) *Response
	// Parse parses the provided JSON string.
	Parse(parse string) gjson.Result
	// Valid checks if the provided JSON string is valid.
	Valid(parse string) bool
	// ToJSON converts the provided object to JSON.
	ToJSON(obj interface{}) string
	// FlattenJSON flattens the provided JSON.
	FlattenJSON(jsonStr string) string
	// Get returns the value for a given JSON path.
	Get(jsonStr, path string) gjson.Result
	// GetArrayElements returns all elements of a JSON array given by the path.
	GetArrayElements(jsonStr, path string) []string
	// GetMapKeyValuePairs returns all key-value pairs from a JSON object at the specified path.
	GetMapKeyValuePairs(jsonStr, path string) map[string]string

	ParseCommandResponse(parse interface{}) *runtime.CommandResponse
}

func New() JSON {
	return &JSONUtils{}
}

// JSONUtils struct contains no fields.
type JSONUtils struct{}

// Response holds the comparison results.
type Response struct {
	NonMatchingFields map[string]interface{}
	MatchingFields    map[string]interface{}
}

// getAllFields recursively retrieves all unique field paths from JSON.
func getAllFields(jsonStr string, prefix string) []string {
	var result []string
	gjson.Parse(jsonStr).ForEach(func(key, value gjson.Result) bool {
		qualifiedKey := key.String()
		if prefix != "" {
			qualifiedKey = prefix + "." + key.String()
		}
		if value.IsObject() || value.IsArray() {
			// Recurse into nested objects and arrays.
			result = append(result, getAllFields(value.Raw, qualifiedKey)...)
		} else {
			// Add leaf node.
			result = append(result, qualifiedKey)
		}
		return true // Keep iterating
	})
	return result
}

// Compare compares two JSON strings based on the provided fields and options.
func (jc *JSONUtils) Compare(j1, j2 string, fields []string, fieldsAsException []string, returnMatching, returnNonMatching bool) *Response {
	result := &Response{
		NonMatchingFields: make(map[string]interface{}),
		MatchingFields:    make(map[string]interface{}),
	}

	if fields == nil {
		// If fields is nil, dynamically determine the fields from both JSON strings.
		fieldsMap := make(map[string]bool)
		for _, field := range getAllFields(j1, "") {
			fieldsMap[field] = true
		}
		for _, field := range getAllFields(j2, "") {
			fieldsMap[field] = true
		}
		// Convert the map keys to a slice.
		for field := range fieldsMap {
			fields = append(fields, field)
		}
	}

	// Helper function to check if a field is in the list of exceptions.
	isException := func(field string) bool {
		for _, f := range fieldsAsException {
			if f == field {
				return true
			}
		}
		return false
	}

	// Iterate over the fields to compare.
	for _, field := range fields {
		if isException(field) {
			continue // Skip fields marked as exceptions.
		}

		val1 := gjson.Get(j1, field)
		val2 := gjson.Get(j2, field)

		// Compare both JSON values.
		if val1.String() == val2.String() {
			if returnMatching {
				result.MatchingFields[field] = val1.Value()
			}
		} else {
			if returnNonMatching {
				result.NonMatchingFields[field] = map[string]interface{}{
					"json1": val1.Value(),
					"json2": val2.Value(),
				}
			}
		}
	}

	return result
}

func (jc *JSONUtils) Parse(parse string) gjson.Result {
	return gjson.Parse(parse)
}

func (jc *JSONUtils) Valid(parse string) bool {
	p := gjson.Valid(parse)
	return p
}

func (jc *JSONUtils) ParseCommandResponse(parse interface{}) *runtime.CommandResponse {
	p := &runtime.CommandResponse{}
	if str, ok := parse.(string); ok {
		err := json.Unmarshal([]byte(str), &p)
		if err != nil {
			return nil
		}
		return p
	}
	if bytesData, ok := parse.([]byte); ok {
		err := json.Unmarshal(bytesData, &p)
		if err != nil {
			return nil
		}
		return p
	}
	return nil
}

func (jc *JSONUtils) ToJSON(obj interface{}) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// FlattenJSON flattens the provided JSON
func (jc *JSONUtils) FlattenJSON(jsonStr string) string {
	j, err := FlattenJSON(jsonStr, DotSeparator)
	if err != nil {
		return ""
	}
	return j
}

// Get returns the value for a given JSON path.
func (jc *JSONUtils) Get(jsonStr, path string) gjson.Result {
	return gjson.Get(jsonStr, path)
}

// GetArrayElements returns all elements of a JSON array given by the path.
func (jc *JSONUtils) GetArrayElements(jsonStr, path string) []string {
	results := gjson.Get(jsonStr, path).Array()
	elements := make([]string, len(results))
	for i, result := range results {
		elements[i] = result.String()
	}
	return elements
}

// GetMapKeyValuePairs returns all key-value pairs from a JSON object at the specified path.
func (jc *JSONUtils) GetMapKeyValuePairs(jsonStr, path string) map[string]string {
	result := gjson.Get(jsonStr, path).Map()
	keyValuePairs := make(map[string]string)
	for key, value := range result {
		keyValuePairs[key] = value.String()
	}
	return keyValuePairs
}

// ------------flattener-----------

// Separator ...
type Separator string

func (s Separator) String() string {
	return string(s)
}

// DotSeparator Separators ...
const (
	DotSeparator Separator = "."
)

type flattenerConfig struct {
	ignoreArray bool
	depth       *int
	prefixes    map[string]bool
}

// Option ...
type Option func(f *flattenerConfig)

// IgnoreArray option, if enabled, ignores arrays while flattening
func IgnoreArray() Option {
	return func(f *flattenerConfig) {
		f.ignoreArray = true
	}
}

// WithDepth option, if provided, limits the flattening to the specified depth
func WithDepth(depth int) Option {
	return func(f *flattenerConfig) {
		f.depth = &depth
	}
}

// FlattenJSON flattens the provided JSON
// The flattening can be customised by providing flattening Options
func FlattenJSON(JSONStr string, separator Separator, options ...Option) (string, error) {
	data := make(map[string]interface{})
	finalMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(JSONStr), &data); err != nil {
		return "", err
	}

	config := &flattenerConfig{}
	for _, option := range options {
		option(config)
	}

	if err := flatten(data, "", separator, config, finalMap, 0); err != nil {
		return "", err
	}

	return mustToJSONStr(finalMap), nil
}

// flatten ....
func flatten(data interface{}, prefix string, separator Separator, config *flattenerConfig, finalMap map[string]interface{}, depth int) error {

	if config.depth != nil && depth == *config.depth {
		finalMap[prefix] = data
		return nil
	}

	switch data.(type) {
	case map[string]interface{}:
		for key, val := range data.(map[string]interface{}) {
			if err := flatten(val, appendToPrefix(prefix, key, separator), separator, config, finalMap, depth+1); err != nil {
				return err
			}
		}
	case []interface{}:
		if config.ignoreArray {
			finalMap[prefix] = data
			return nil
		}
		for index, val := range data.([]interface{}) {
			if err := flatten(val, appendToPrefix(prefix, fmt.Sprintf("%v", index), separator), separator, config, finalMap, depth+1); err != nil {
				return err
			}
		}

	default:
		finalMap[prefix] = data
	}

	return nil
}

func appendToPrefix(prefix string, key string, separator Separator) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v%v%v", prefix, separator.String(), key)
}

func mustToJSONStr(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
