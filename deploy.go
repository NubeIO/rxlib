package rxlib

import (
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type Deploy struct {
	ObjectDeploy `json:"objectDeploy"`
	Timeout      int `json:"timeout"`
}

type ObjectDeploy struct {
	Deleted []string                `json:"deleted"`
	New     []*runtime.ObjectConfig `json:"new"`
	Updated []*runtime.ObjectConfig `json:"updated"`
}

func SerializeCurrentFlowArray(objects []Object) []*runtime.ObjectConfig {
	var serializedObjects []*runtime.ObjectConfig
	for _, object := range objects {
		serializedObjects = append(serializedObjects, serializeCurrentFlowArray(object))
	}
	return serializedObjects
}

func serializeCurrentFlowArray(object Object) *runtime.ObjectConfig {
	if object == nil {
		return nil
	}
	meta := object.GetMeta()
	if meta == nil {
		meta = &runtime.Meta{
			Position: &runtime.Position{
				PositionY: 0,
				PositionX: 0,
			},
		}
	}
	objectConfig := &runtime.ObjectConfig{
		Id:          object.GetID(),
		Info:        object.GetInfo(),
		Inputs:      PortsToProto(object.GetInputs()),
		Outputs:     PortsToProto(object.GetOutputs()),
		Connections: object.GetConnections(),
		Settings:    object.GetSettings(),
		Stats:       object.GetStats(),
		Meta:        meta,
	}
	return objectConfig
}

func ProtosToPort(obj []*runtime.Port) []*Port {
	var out []*Port
	for _, port := range obj {
		out = append(out, ProtoToPort(port))
	}
	return out
}

func PortsToProto(obj []*Port) []*runtime.Port {
	var out []*runtime.Port
	for _, port := range obj {
		out = append(out, PortToProto(port))
	}
	return out
}

func PortToProto(obj *Port) *runtime.Port {
	return &runtime.Port{
		Id:              obj.ID,
		Name:            obj.Name,
		PortUUID:        obj.UUID,
		Direction:       string(obj.Direction),
		DataType:        string(obj.DataType),
		DefaultPosition: int32(obj.DefaultPosition),
	}
}

func ProtoToPort(obj *runtime.Port) *Port {
	return &Port{
		ID:              obj.Id,
		Name:            obj.Name,
		UUID:            obj.PortUUID,
		Direction:       PortDirection(obj.Direction),
		DataType:        priority.Type(obj.DataType),
		DefaultPosition: int(obj.DefaultPosition),
	}
}
