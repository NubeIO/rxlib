package extensionlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
)

func (inst *Extensions) AddPallet(name string, ext interface{}) {
	baseObj := reactive.New(name, nil)
	extInstance, ok := ext.(Payload)
	if !ok {
		print(ok)
	}
	obj := extInstance.New(baseObj)
	obj.GetInfo().PluginName = inst.name
	inst.pallet = append(inst.pallet, obj)
}
