package runtimebase

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

func ObjectConfigToProto(obj *rxlib.ObjectConfig) *runtime.Object {
	return &runtime.Object{
		Id: obj.ID,
		//Info:        ObjectInfoToProto(obj.Info),
		Inputs:      nil,
		Outputs:     nil,
		Meta:        nil,
		Stats:       nil,
		Connections: nil,
	}
}

func ObjectConfigFromProto(protoObj *runtime.Object) *rxlib.ObjectConfig {
	return &rxlib.ObjectConfig{
		ID: protoObj.Id,
	}
}

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
	cmd := convertCommandResponse(c)
	if len(c.CommandResponse) > 0 {
		var out []*runtime.CommandResponse
		for _, response := range c.CommandResponse {
			out = append(out, convertCommandResponse(response))
		}
		cmd.Response = out
	}
	if len(c.SerializeObjects) > 0 {
		var out []*runtime.Object
		for _, response := range c.SerializeObjects {
			out = append(out, ObjectConfigToProto(response))
		}
		cmd.SerializeObjects = out
	}
	return cmd
}
