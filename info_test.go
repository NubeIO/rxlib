package rxlib

import (
	"fmt"
	"testing"
)

func TestNewObjectInfo(t *testing.T) {
	builder := NewObjectInfo()

	// Set some values using the InfoBuilder methods
	info := builder.
		SetID("123").
		SetPluginName("MyPlugin")

	fmt.Printf("%+v\n", info.Build())

}
