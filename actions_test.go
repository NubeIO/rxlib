package rxlib

import (
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestActions(t *testing.T) {

	type responseModel struct {
		Value int `json:"value"`
	}
	m, _ := ReflectStructFields(responseModel{})

	actionList := NewActionList()
	actionList.AddGetAction(&ActionBody{
		Name:          "get user",
		Description:   "get a user",
		Path:          "user",
		Method:        "GET",
		BodyModel:     nil,
		ResponseModel: m,
	})

	// Get all ActionLists from the ActionList interface
	actionListsMap := actionList.All()

	pprint.PrintJSON(actionListsMap)

}
