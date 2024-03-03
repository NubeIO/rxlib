package runtimebase

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

//
//func ProtosToPort(obj []*runtime.Port) []*rxlib.Port {
//	var out []*rxlib.Port
//	for _, port := range obj {
//		out = append(out, ProtoToPort(port))
//	}
//	return out
//}
//
//func PortsToProto(obj []*rxlib.Port) []*runtime.Port {
//	var out []*runtime.Port
//	for _, port := range obj {
//		out = append(out, PortToProto(port))
//	}
//	return out
//}
//
//func PortToProto(obj *rxlib.Port) *runtime.Port {
//	return &runtime.Port{
//		Id:              obj.ID,
//		Name:            obj.Name,
//		PortUUID:        obj.UUID,
//		Direction:       string(obj.Direction),
//		DataType:        string(obj.DataType),
//		DefaultPosition: int32(obj.DefaultPosition),
//	}
//}
//
//func ProtoToPort(obj *runtime.Port) *rxlib.Port {
//	return &rxlib.Port{
//		ID:              obj.Id,
//		Name:            obj.Name,
//		UUID:            obj.PortUUID,
//		Direction:       rxlib.PortDirection(obj.Direction),
//		DataType:        priority.Type(obj.DataType),
//		DefaultPosition: int(obj.DefaultPosition),
//	}
//}

//func ObjectConfigToProto(obj *runtime.ObjectConfig) *runtime.ObjectConfig {
//	return &runtime.ObjectConfig{
//		Id: obj.Id,
//		//Info:        ObjectInfoToProto(obj.Info),
//		Inputs:      nil,
//		Outputs:     nil,
//		Meta:        nil,
//		Stats:       nil,
//		Connections: nil,
//	}
//}

/*
	func ObjectConfigFromProto(protoObj *runtime.ObjectConfig) *runtime.ObjectConfig {
		return &rxlib.ObjectConfig{
			ID: protoObj.Id,
		}
	}
*/
func ConvertCommand(c *runtime.Command) *rxlib.ExtendedCommand {
	//out := &rxlib.ExtendedCommand{
	//	TargetGlobalID:   c.GetTargetGlobalID(),
	//	SenderGlobalID:   c.GetSenderGlobalID(),
	//	SenderObjectUUID: c.GetSenderObjectUUID(),
	//	TransactionUUID:  c.GetTransactionUUID(),
	//	Key:              c.GetKey(),
	//	Args:             c.GetArgs(),
	//	Data:             c.GetData(),
	//	Body:             c.Body.GetValue(),
	//}
	return nil
}

func convertCommandResponse(c *rxlib.CommandResponse) *runtime.CommandResponse {
	out := &runtime.CommandResponse{
		SenderID:   c.SenderID,
		Count:      int32(nils.GetInt(c.Count)),
		MapStrings: c.MapStrings,
		Number:     nils.GetFloat64(c.Float),
		Boolean:    nils.GetBool(c.Bool),
		Error:      c.Error,
		ReturnType: c.ReturnType,
		Any:        c.Any,
	}
	return out
}

func ConvertCommandResponse(c *rxlib.CommandResponse) *runtime.CommandResponse {
	//cmd := convertCommandResponse(c)
	//if len(c.CommandResponse) > 0 {
	//	var out []*runtime.CommandResponse
	//	for _, response := range c.CommandResponse {
	//		out = append(out, convertCommandResponse(response))
	//	}
	//	cmd.Response = out
	//}
	//if len(c.SerializeObjects) > 0 {
	//	var out []*runtime.ObjectConfig
	//	for _, response := range c.SerializeObjects {
	//		//out = append(out, ObjectConfigToProto(response))
	//	}
	//	cmd.SerializeObjects = out
	//}
	//return cmd
	return nil
}
