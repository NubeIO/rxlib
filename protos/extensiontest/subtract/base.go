package subtract

import (
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/pluginlib"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"log"
	"os"
	"time"
)

type Instance struct {
	reactive.Object
	locked        bool
	lastTrigger   time.Time
	in1           float64
	in2           float64
	outputUpdated func(message *runtime.Command)
	portOne       float64
	portTwo       float64
	lastValue     float64
	hasPublished  bool
}

var infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func New(outputUpdated func(message *runtime.Command)) pluginlib.PluginObject {
	obj := new(Instance)
	obj.outputUpdated = outputUpdated
	return obj
}

func (inst *Instance) New(object reactive.Object, opts ...any) reactive.Object {
	info := rxlib.NewObjectInfo().
		SetID("subtract").
		SetPluginName("ext-math").
		SetCategory("math").
		SetCallResetOnDeploy().
		SetObjectType(rxlib.Logic).
		SetAllPermissions().
		Build()

	object.SetInfo(info)
	object.NewOutputPort(&runtime.Port{
		Id:        "output",
		Name:      "output",
		Direction: string(rxlib.Output),
		DataType:  priority.TypeFloat,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "input-1",
		Name:      "input-1",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeFloat,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "input-2",
		Name:      "input-2",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeFloat,
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
	infoLog.Println("subtract Handler")
	if p == nil {
		return
	}
	cmd := p.GetCommand()
	if cmd == nil {
		return
	}

	for _, value := range cmd.GetPortValues() {
		for _, d := range value.PortIDs {
			if d == "input-1" {
				inst.portOne = *value.FloatValue
			}
			if d == "input-2" {
				inst.portTwo = *value.FloatValue
			}
		}
	}
	inst.publishOutput()

}

func (inst *Instance) publishOutput() {
	v := inst.portOne - inst.portTwo
	var cov bool
	if v != inst.lastValue {
		cov = true
	}
	if cov || !inst.hasPublished {
		fmt.Println("SUBTRACT", v, inst.GetMeta().ObjectUUID)
		inst.OutputUpdated(&runtime.Command{
			Key:              "update-outputs",
			TargetObjectUUID: inst.GetMeta().GetObjectUUID(),
			PortValues: []*runtime.PortValue{&runtime.PortValue{
				PortID:     "output",
				FloatValue: &v,
				DataType:   priority.TypeFloat,
			}},
		})
		inst.hasPublished = true // this is for to make sure we publish the first value
	} else {
		//fmt.Println("ADD SKIP", v, inst.GetMeta().ObjectUUID)

	}
	inst.lastValue = v

}
