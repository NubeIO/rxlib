package rxlib

import (
	"fmt"
	"reflect"
)

type ActionList interface {
	All() *ActionLists // Change the return type to map[string]*ActionLists
	AddGetAction(action *ActionBody)
	AddSetAction(action *ActionBody)
	GetGetActions() []*ActionBody
	GetSetActions() []*ActionBody
}

func NewActionList() ActionList {
	return &ActionLists{}
}

func (al *ActionLists) All() *ActionLists {
	return &ActionLists{
		GetActions: al.GetGetActions(),
		SetActions: al.GetSetActions(),
	}
}

// AddGetAction adds a GetAction to the list of GetActions.
func (al *ActionLists) AddGetAction(action *ActionBody) {
	al.GetActions = append(al.GetActions, action)
}

// AddSetAction adds a SetAction to the list of SetActions.
func (al *ActionLists) AddSetAction(action *ActionBody) {
	al.SetActions = append(al.SetActions, action)
}

// GetGetActions returns the list of GetActions.
func (al *ActionLists) GetGetActions() []*ActionBody {
	return al.GetActions
}

// GetSetActions returns the list of SetActions.
func (al *ActionLists) GetSetActions() []*ActionBody {
	return al.SetActions
}

type ActionLists struct {
	GetActions []*ActionBody `json:"getActions,omitempty"`
	SetActions []*ActionBody `json:"setActions,omitempty"`
}

type ActionResponse struct {
	Body    any    `json:"body,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ActionBody struct {
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
