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

const (
	null                 = "null"
	uid                  = "uuid"
	name                 = "name"
	id                   = "id"
	objType              = "type"
	parent               = "parent"
	inputPort            = "input"
	outputPort           = "output"
	propertyID           = "id"
	propertyName         = "name"
	propertyValue        = "value"
	propertyUnitFrom     = "unitFrom"
	propertyUnitTo       = "unitTo"
	propertyDisplayValue = "displayValue"
	propertyRawValue     = "rawValue"
	propertyType         = "type"
)

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

	rawValue := handleQuery(objects, obj, queryParts)
	return formatOutput(rawValue, formatString, queryTerm)
}

func handleQuery(objects []Object, obj Object, queryParts []string) interface{} {
	switch queryParts[0] {
	case uid, name, id, objType:
		return extractMetaData(obj, queryParts[0])
	case parent:
		return handleParent(objects, obj, queryParts)
	case inputPort, outputPort:
		return handlePorts(obj, queryParts)
	default:
		return fmt.Sprintf("invalid-query: (%s)", strings.Join(queryParts, "."))
	}
}

func extractMetaData(obj Object, property string) interface{} {
	switch property {
	case uid:
		return obj.GetMeta().GetObjectUUID()
	case name:
		return obj.GetMeta().GetObjectName()
	case id:
		return obj.GetID()
	case objType:
		return string(obj.GetObjectType())
	default:
		return nil
	}
}

func handleParent(objects []Object, obj Object, queryParts []string) interface{} {
	p := GetByUUID(objects, obj.GetMeta().GetParentUUID())
	if p != nil {
		return extractInfo(objects, p, strings.Join(queryParts[1:], "."))
	}
	return nil
}

func handlePorts(obj Object, queryParts []string) interface{} {
	if len(queryParts) < 3 {
		return fmt.Sprintf("invalid-query: (%s)", strings.Join(queryParts, "."))
	}

	portType := queryParts[0]
	portID := queryParts[1]
	property := queryParts[2]

	var port *Port
	if portType == inputPort {
		port = obj.GetInput(portID)
	} else {
		port = obj.GetOutput(portID)
	}

	if port == nil {
		return fmt.Sprintf("port-not-found: (%s)", strings.Join(queryParts, "."))
	}

	return extractPortProperty(port, property)
}

func extractPortProperty(port *Port, property string) interface{} {
	switch property {
	case propertyID:
		return port.GetID()
	case propertyName:
		return port.GetName()
	case propertyValue:
		value, isNil := port.GetPayloadValue()
		if isNil {
			return null
		} else {
			return value
		}
	case propertyUnitFrom:
		value := port.GetTransformationUnitFrom()
		if value == "" {
			return null
		} else {
			return value
		}
	case propertyUnitTo:
		value := port.GetTransformationUnitTo()
		if value == "" {
			return null
		} else {
			return value
		}
	case propertyDisplayValue:
		value, isNil := port.GetTransformationDisplayValue()
		if isNil {
			return null
		} else {
			return value
		}
	case propertyRawValue:
		value, isNil := port.GetTransformationExistingValueFloat()
		if isNil {
			return null
		} else {
			return value
		}
	case propertyType:
		return string(port.GetDataType())
	default:
		return fmt.Sprintf("invalid-property: (%s)", property)
	}
}

func formatOutput(rawValue interface{}, formatString string, queryTerm string) string {
	if formatString != "" {
		switch v := rawValue.(type) {
		case string:
			if numValue, err := strconv.ParseFloat(v, 64); err == nil {
				return fmt.Sprintf(formatString, numValue)
			} else if boolValue, err := strconv.ParseBool(v); err == nil {
				return fmt.Sprintf(formatString, boolValue)
			}
		case float64, float32, int, int64, int32:
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
