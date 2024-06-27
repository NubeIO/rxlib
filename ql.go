package rxlib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetByUUID(objects []Object, uuid string) Object {
	for _, obj := range objects {
		if obj.GetMeta().ObjectUUID == uuid {
			return obj
		}
	}
	return nil
}

func extractInfo(objects []Object, obj Object, term string) string {
	parts := strings.Split(term, "|")
	queryTerm := parts[0]
	formatString := ""
	if len(parts) > 1 {
		formatString = parts[1]
	}

	queryParts := strings.Split(queryTerm, ".")
	if len(queryParts) == 0 {
		return fmt.Sprintf("invalid-query: (%s)", queryTerm)
	}

	var rawValue any
	switch queryParts[0] {
	case "uuid":
		if uuid := obj.GetMeta().GetObjectUUID(); uuid != "" {
			rawValue = uuid
		}
	case "name":
		if name := obj.GetMeta().GetObjectName(); name != "" {
			rawValue = name
		}
	case "id":
		if id := obj.GetID(); id != "" {
			rawValue = id
		}
	case "type":
		if objType := string(obj.GetObjectType()); objType != "" {
			rawValue = objType
		}
	case "parent":
		parent := GetByUUID(objects, obj.GetMeta().GetParentUUID())
		if parent != nil {
			rawValue = extractInfo(objects, parent, strings.Join(queryParts[1:], "."))
		}
	case "input", "output":
		if len(queryParts) < 3 {
			return fmt.Sprintf("invalid-query: (%s)", queryTerm)
		}
		portType := queryParts[0]
		portID := queryParts[1]
		property := queryParts[2]

		var port *Port
		if portType == "input" {
			port = obj.GetInput(portID)
		} else {
			port = obj.GetOutput(portID)
		}

		if port == nil {
			return fmt.Sprintf("port-not-found: (%s)", queryTerm)
		}

		switch property {
		case "id":
			rawValue = port.GetID()
		case "name":
			rawValue = port.GetName()
		case "value":
			value, isNil := port.GetPayloadValue()
			if isNil {
				rawValue = "null"
			} else {
				rawValue = value
			}
		case "type":
			rawValue = string(port.GetDataType())
		default:
			rawValue = fmt.Sprintf("invalid-property: (%s)", queryTerm)
		}
	default:
		rawValue = fmt.Sprintf("invalid-query: (%s)", queryTerm)
	}
	if formatString != "" {
		switch v := rawValue.(type) {
		case string:
			if numValue, err := strconv.ParseFloat(v, 64); err == nil {
				return fmt.Sprintf(formatString, numValue)
			} else if boolValue, err := strconv.ParseBool(v); err == nil {
				return fmt.Sprintf(formatString, boolValue)
			}
		case float64, float32:
			return fmt.Sprintf(formatString, v)
		case int, int64, int32:
			return fmt.Sprintf(formatString, v)
		case bool:
			return fmt.Sprintf(formatString, v)
		default:
			return fmt.Sprintf("unsupported-type: (%s)", queryTerm)
		}
	}

	return fmt.Sprintf("%v", rawValue)
}

type QueryResult struct {
	QueryString string   `json:"queryString"`
	Results     []string `json:"results"`
	Result      any      `json:"result"`
}

type Query struct {
	QueryString string `json:"queryString"`
}

func ParseTemplates(objects []Object, objs []Object, templates []string) (map[string]*QueryResult, error) {
	if objects == nil {
		return nil, errors.New("objects is nil")
	}
	if objs == nil {
		return nil, errors.New("objs is nil")
	}
	if len(objs) != len(templates) {
		return nil, errors.New("the number of objects and templates must be the same")
	}
	results := make(map[string]*QueryResult)

	for i, obj := range objs {
		template := templates[i]
		result, resultString, err := ParseTemplate(objects, obj, template)
		if err != nil {
			return nil, err
		}
		objKey := fmt.Sprintf("%s-%d", obj.GetMeta().GetObjectUUID(), i)
		results[objKey] = &QueryResult{
			QueryString: template,
			Results:     result,
			Result:      resultString,
		}
	}
	return results, nil
}

func ParseTemplate(objects []Object, obj Object, template string) ([]string, string, error) {
	if objects == nil {
		return nil, "", errors.New("objects is nil")
	}
	if obj == nil {
		return nil, "", errors.New("object is nil")
	}
	var result []string
	var finalStringParts []string
	parts := strings.Split(template, "*")
	for i, part := range parts {
		if i%2 == 0 {
			if part != "" {
				finalStringParts = append(finalStringParts, part)
			}
		} else {
			extracted := extractInfo(objects, obj, strings.TrimSpace(part))
			result = append(result, extracted)
			finalStringParts = append(finalStringParts, extracted)
		}
	}

	finalString := strings.Join(finalStringParts, "")
	return result, finalString, nil
}
