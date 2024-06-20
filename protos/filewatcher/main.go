package main

import (
	"github.com/NubeIO/rxlib/protos/extensionlib"
	"github.com/NubeIO/rxlib/protos/filewatcher/jsonpath"
	"github.com/NubeIO/rxlib/protos/filewatcher/watcher"
)

func main() {
	factory := extensionlib.New("test")
	factory.AddPallet("watcher", watcher.New)
	factory.AddPallet("jsonpath", jsonpath.New)
	factory.Register()
}
