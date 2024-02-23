package rxlib

import "github.com/NubeIO/rxlib/plugins"

type PalletTree struct {
	Drivers  []*ObjectConfig `json:"drivers"`
	Logic    []*ObjectConfig `json:"logic"`
	Services []*ObjectConfig `json:"services"`
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
