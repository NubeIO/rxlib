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

	actionList := NewInvokeList()
	actionList.AddAction(&Invoke{
		Name:          "get user",
		Description:   "get a user",
		Path:          "user",
		Method:        "GET",
		BodyModel:     nil,
		ResponseModel: m,
	})

	// Get all InvokeList from the InvokeBuilder interface
	actionListsMap := actionList.All()

	pprint.PrintJSON(actionListsMap)

}
