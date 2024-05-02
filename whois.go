package rxlib

import (
	"github.com/NubeIO/rxlib/helpers"
	systeminfo "github.com/NubeIO/rxlib/libs/system"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type Iam struct {
	GlobalID string           `json:"globalID"`
	DeviceID int              `json:"deviceID"`
	System   *systeminfo.Info `json:"system"`
	Error    string           `json:"error,omitempty"`
}

func (inst *RuntimeImpl) whois(parsedArgs *ParsedCommand) Object {
	var start int
	var finish int
	if parsedArgs.GetGlobal() {
		iam := inst.Iam(start, finish)
		return iam
	}
	start = parsedArgs.GetStart()
	finish = parsedArgs.GetFinish()
	iam := inst.Iam(start, finish)
	return iam
}

func (inst *RuntimeImpl) IamConfig(rangeStart, rangeEnd int) *runtime.ObjectConfig {
	obj := inst.GetFirstByID("rubix-manager")
	if obj == nil {
		return nil
	}
	obj.Invoke(nil) // update all the info
	globalID := obj.GetFlag("globalID")
	id, err := helpers.ProcessID(globalID)
	if err != nil {
		return nil
	}
	if rangeStart == 0 && rangeEnd == 0 {
		return inst.serializeObject(false, obj)
	}
	if isBetween(id, rangeStart, rangeEnd) {
		return inst.serializeObject(false, obj)
	}
	return nil
}

func (inst *RuntimeImpl) Iam(rangeStart, rangeEnd int) Object {
	obj := inst.GetFirstByID("rubix-manager")
	if obj == nil {
		return nil
	}
	obj.Invoke(nil) // update all the info
	globalID := obj.GetFlag("globalID")
	id, err := helpers.ProcessID(globalID)
	if err != nil {
		return nil
	}
	if rangeStart == 0 && rangeEnd == 0 {
		return obj
	}
	if isBetween(id, rangeStart, rangeEnd) {
		return obj
	}
	return nil
}

func isBetween(a, b, c int) bool {
	return a >= b && a <= c
}
