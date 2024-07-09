package main

import (
	"github.com/NubeIO/rxlib/libs/pluginlib"
	"github.com/NubeIO/rxlib/protos/filewatcher/jsonpath"
	"github.com/NubeIO/rxlib/protos/filewatcher/watcher"
)

func main() {
	factory := pluginlib.New("test")
	factory.AddPallet("watcher", watcher.New)
	factory.AddPallet("jsonpath", jsonpath.New)
	factory.Register("test")
}
