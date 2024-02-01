package rxlib

import (
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestNewExtensionBuilder(t *testing.T) {
	extensions := NewExtensionBuilder().
		NewExtension().
		WithExtensionName("MyExtension1").
		WithFromPlugin("Plugin1").
		WithParentObjectUUID("12345").
		AddExtensionAutoConnect("output1", "input1").
		NewExtension().
		WithExtensionName("MyExtension2").
		WithFromPlugin("Plugin2").
		WithParentObjectUUID("67890").
		AddExtensionAutoConnect("output2", "input2").
		Build()

	pprint.PrintJSON(extensions)
}
