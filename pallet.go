package rxlib

import (
	"github.com/NubeIO/rxlib/plugins"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type PalletTree struct {
	Drivers  []*runtime.ObjectConfig `json:"drivers"`
	Logic    []*runtime.ObjectConfig `json:"logic"`
	Services []*runtime.ObjectConfig `json:"services"`
}

func (inst *RuntimeImpl) AllPlugins() []*plugins.Export {
	return inst.PluginsExport
}

func (inst *RuntimeImpl) PluginTree() map[string][]Object {
	out := make(map[string][]Object)
	//for i, object := range inst.pluginObjects {
	//
	//}
	return out
}

func (inst *RuntimeImpl) GetObjectsPallet() *PalletTree {
	//for i, object := range inst.PluginObjects {
	//
	//}
	return &PalletTree{}
}
