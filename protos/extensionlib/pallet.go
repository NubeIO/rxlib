package extensionlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
)

func (inst *Extensions) AddPallet(name string, generate GeneratePlugin) {
	baseObj := reactive.New(name, nil)
	ext := generate(nil)
	obj := ext.New(baseObj)
	obj.GetInfo().PluginName = inst.name
	inst.pallet[name] = obj
	inst.registry[name] = generate
}
