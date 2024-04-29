package jsonutils

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func (jc *JSONUtils) Parse(parse string) gjson.Result {
	return gjson.Parse(parse)
}

func (jc *JSONUtils) Valid(parse string) bool {
	p := gjson.Valid(parse)
	return p
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

// GetValue returns the value for a given JSON path.
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
