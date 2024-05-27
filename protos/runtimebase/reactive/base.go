package reactive

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

func ConvertObjects(objects []Object) []*runtime.ObjectConfig {
	var out []*runtime.ObjectConfig
	for _, object := range objects {
		out = append(out, Convert(object))
	}
	return out
}

func ConvertObject(object *runtime.ObjectConfig) Object {
	baseObj := &BaseObject{
		Inputs:  object.GetInputs(),
		Outputs: object.GetOutputs(),
		Info:    object.GetInfo(),
		Meta:    object.GetMeta(),
	}
	return baseObj
}

func ConvertObjectConfig(object *runtime.ObjectConfig) *BaseObject {
	baseObj := &BaseObject{
		Inputs:  object.GetInputs(),
		Outputs: object.GetOutputs(),
		Info:    object.GetInfo(),
		Meta:    object.GetMeta(),
	}
	return baseObj
}

func Convert(object Object) *runtime.ObjectConfig {
	return &runtime.ObjectConfig{
		Id:          object.GetInfo().ObjectID,
		Info:        object.GetInfo(),
		Inputs:      object.GetInputs(),
		Outputs:     object.GetOutputs(),
		Meta:        object.GetMeta(),
		Stats:       nil,
		Connections: nil,
		Settings:    nil,
	}
}

type Object interface {
	Handler(payload *runtime.MessageRequest)
	NewOutputPort(port *runtime.Port) error
	NewInputPort(port *runtime.Port) error
	SetInfo(info *runtime.Info)

	GetOutputs() []*runtime.Port
	GetInputs() []*runtime.Port
	GetInfo() *runtime.Info
	GetMeta() *runtime.Meta
}

func New(id string, obj Object) Object {
	if obj != nil {
		return obj
	}
	baseObj := &BaseObject{
		id:      id,
		Inputs:  []*runtime.Port{},
		Outputs: []*runtime.Port{},
		Info: &runtime.Info{
			ObjectID: id,
		},
		Meta: &runtime.Meta{
			ObjectUUID: "",
			ObjectName: id,
		},
	}
	return baseObj
}

type BaseObject struct {
	id string
	// object Inputs
	Inputs []*runtime.Port

	// object Outputs
	Outputs []*runtime.Port

	// object Info like uuid, id
	Info *runtime.Info

	// Meta is data sent from the UI like name, object position
	Meta *runtime.Meta
}

func (inst *BaseObject) GetOutputs() []*runtime.Port {
	return inst.Outputs
}

func (inst *BaseObject) GetInputs() []*runtime.Port {
	return inst.Inputs
}

func (inst *BaseObject) GetInfo() *runtime.Info {
	return inst.Info
}
func (inst *BaseObject) GetMeta() *runtime.Meta {
	return inst.Meta
}

func (inst *BaseObject) Handler(payload *runtime.MessageRequest) {}

func (inst *BaseObject) NewOutputPort(port *runtime.Port) error {
	return inst.newPort(port, "output")
}

func (inst *BaseObject) NewInputPort(port *runtime.Port) error {
	return inst.newPort(port, "input")
}

func (inst *BaseObject) newPort(port *runtime.Port, direction string) error {
	if port.PortUUID == "" {
		port.PortUUID = helpers.UUID()
	}
	if direction == string(rxlib.Input) {
		inst.Inputs = append(inst.Inputs, port)
	} else if direction == string(rxlib.Output) {
		inst.Outputs = append(inst.Outputs, port)
	}
	return nil
}

func (inst *BaseObject) SetInfo(info *runtime.Info) {
	inst.Info = info
	if inst.Info.Requirements == nil {
		inst.Info.Requirements = &runtime.Requirements{}
	}
}
