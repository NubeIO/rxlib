package main

import (
	"github.com/NubeIO/rxlib/protos/extensionlib"
	"github.com/NubeIO/rxlib/protos/extensiontest/add"
	"github.com/NubeIO/rxlib/protos/extensiontest/subtract"
)

func main() {
	factory := extensionlib.New("test")
	factory.AddPallet("add", add.New)
	factory.AddPallet("subtract", subtract.New)
	factory.Register()
}
