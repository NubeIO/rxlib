package counter

import (
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"time"
)

type Instance struct {
	reactive.Object
	locked      bool
	lastTrigger time.Time
	count       int
}

func New() *Instance {
	obj := new(Instance)
	return obj
}

func (inst *Instance) New(object reactive.Object, opts ...any) reactive.Object {

	info := rxlib.NewObjectInfo().
		SetID("counter-2").
		SetPluginName("1234").
		SetCategory("time").
		SetCallResetOnDeploy().
		SetObjectType(rxlib.Logic).
		SetAllPermissions().
		Build()

	object.SetInfo(info)
	err := object.NewOutputPort(&runtime.Port{
		Id:        "output",
		Name:      "output",
		Direction: string(rxlib.Output),
		DataType:  "float",
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	object.NewInputPort(&runtime.Port{
		Id:        "input",
		Name:      "input",
		Direction: string(rxlib.Input),
		DataType:  "float",
	})
	return &Instance{
		Object: object,
	}
}

func (inst *Instance) Start() error {
	return nil

}

func (inst *Instance) Reset() error {
	return nil
}

func (inst *Instance) Handler(payload *runtime.MessageRequest) {
	fmt.Println(payload)
	fmt.Println(11111)
}
