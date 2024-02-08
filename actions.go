package rxlib

import (
	"fmt"
	"reflect"
)

type InvokeBuilder interface {
	All() []*Invoke
	AddAction(action *Invoke)
}

func NewInvokeList() InvokeBuilder {
	return &InvokeList{}
}

func (al *InvokeList) All() []*Invoke {
	return al.List
}

func (al *InvokeList) AddAction(action *Invoke) {
	al.List = append(al.List, action)
}

type InvokeList struct {
	List []*Invoke `json:"actions,omitempty"`
}

type Invoke struct {
	ObjectUUID    string `json:"objectUUID,omitempty"`
	Name          string `json:"name,omitempty"`        // user
	Description   string `json:"description,omitempty"` // get user
	Example       string `json:"example,omitempty"`     // "RQL.ObjectInvoke"
	Path          string `json:"path,omitempty"`        // user
	Method        string `json:"method,omitempty"`      // GET output
	Body          any    `json:"body,omitempty"`
	BodyModel     any    `json:"bodyModel,omitempty"`     // we may need to send down some JSON {uuid:string}, so we need to send the exacted body to be sent
	ResponseModel any    `json:"responseModel,omitempty"` // {name:string, age:int}
}

func ReflectStructFields(s interface{}) (map[string]string, error) {
	result := make(map[string]string)

	// Use reflection to inspect the struct type
	structType := reflect.TypeOf(s)

	// Check if the provided value is a struct
	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	// Iterate through the struct's fields
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name
		fieldType := field.Type.Name()

		// Add the field name and type to the result map
		result[fieldName] = fieldType
	}

	return result, nil
}
