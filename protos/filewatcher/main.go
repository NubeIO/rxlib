package main

import (
	"github.com/NubeIO/rxlib/protos/filewatcher/jsonpath"
	"github.com/NubeIO/rxlib/protos/filewatcher/watcher"
	"github.com/NubeIO/rxlib/protos/pluginlib"
)

func main() {
	factory := pluginlib.New("test")
	factory.AddPallet("watcher", watcher.New)
	factory.AddPallet("jsonpath", jsonpath.New)
	factory.Register()
}
