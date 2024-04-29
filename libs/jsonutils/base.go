package jsonutils

import (
	"github.com/tidwall/gjson"
)

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
