package jsonpath

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/pluginlib"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/tidwall/gjson"
	"log"
)

type Instance struct {
	reactive.Object
	json          string
	path          string
	output        string
	outputUpdated func(message *runtime.Command)
}

func New(outputUpdated func(message *runtime.Command)) pluginlib.PluginObject {
	obj := new(Instance)
	obj.outputUpdated = outputUpdated
	return obj
}

func (inst *Instance) New(object reactive.Object, opts ...any) reactive.Object {
	info := rxlib.NewObjectInfo().
		SetID("jsonpath").
		SetPluginName("test").
		SetCategory("util").
		SetCallResetOnDeploy().
		SetObjectType(rxlib.Service).
		SetAllPermissions().
		Build()

	object.SetInfo(info)
	object.NewOutputPort(&runtime.Port{
		Id:        "output",
		Name:      "output",
		Direction: string(rxlib.Output),
		DataType:  priority.TypeAny,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "json",
		Name:      "json",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeString,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "path",
		Name:      "path",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeString,
	})
	inst.Object = object
	return inst
}

func (inst *Instance) OutputUpdated(message *runtime.Command) {
	inst.outputUpdated(message)
}

func (inst *Instance) Start() error {
	return nil
}

func (inst *Instance) Reset() error {
	return nil
}

func (inst *Instance) Handler(p *runtime.MessageRequest) {
	log.Println("jsonpath Handler")
	if p == nil {
		return
	}
	cmd := p.GetCommand()
	if cmd == nil {
		return
	}

	for _, value := range cmd.GetPortValues() {
		for _, d := range value.PortIDs {
			if d == "json" {
				inst.json = *value.StringValue
			}
			if d == "path" {
				inst.path = *value.StringValue
			}
		}
	}

	inst.publishOutput()
}

func (inst *Instance) publishOutput() {
	value := gjson.Get(inst.json, inst.path).String()
	inst.OutputUpdated(&runtime.Command{
		Key:              "update-outputs",
		TargetObjectUUID: inst.GetMeta().GetObjectUUID(),
		PortValues: []*runtime.PortValue{&runtime.PortValue{
			PortID:      "output",
			StringValue: &value,
			DataType:    priority.TypeString,
		}},
	})
}
