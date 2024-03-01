package runtime

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/protos/runtime/protoruntime"
)

func ObjectConfigToProto(obj *rxlib.ObjectConfig) *protoruntime.Object {
	return &protoruntime.Object{
		Id: obj.ID,
	}
}

func ObjectConfigFromProto(protoObj *protoruntime.Object) *rxlib.ObjectConfig {
	return &rxlib.ObjectConfig{
		ID: protoObj.Id,
	}
}

func ConvertProtoConnectionToStruct(protoConn *protoruntime.Connection) *rxlib.Connection {
	conn := &rxlib.Connection{
		UUID:                 protoConn.GetUuid(),
		TargetConnectionUUID: protoConn.GetTargetConnectionUUID(),
		SourceUUID:           protoConn.GetSource(),
		SourcePort:           protoConn.GetSourceHandle(),
		SourcePortUUID:       protoConn.GetSourcePortUUID(),
		TargetUUID:           protoConn.GetTarget(),
		TargetPort:           protoConn.GetTargetHandle(),
		TargetPortUUID:       protoConn.GetTargetPortUUID(),
	}

	return conn
}

func ConvertStructConnectionToProto(conn *rxlib.Connection) *protoruntime.Connection {
	protoConn := &protoruntime.Connection{
		Uuid:                 conn.UUID,
		TargetConnectionUUID: conn.TargetConnectionUUID,
		Source:               conn.SourceUUID,
		SourceHandle:         conn.SourcePort,
		SourcePortUUID:       conn.SourcePortUUID,
		Target:               conn.TargetUUID,
		TargetHandle:         conn.TargetPort,
		TargetPortUUID:       conn.TargetPortUUID,
	}

	return protoConn
}

func ConvertCommand(command *protoruntime.CommandRequest) *rxlib.Command {
	c := command.GetCommand()
	out := &rxlib.Command{
		TargetGlobalID:   c.GetTargetGlobalID(),
		SenderGlobalID:   c.GetSenderGlobalID(),
		SenderObjectUUID: c.GetSenderObjectUUID(),
		TransactionUUID:  c.GetTransactionUUID(),
		Key:              c.GetKey(),
		Query:            c.GetQuery(),
		Args:             c.GetArgs(),
		Data:             c.GetData(),
		Body:             c.Body.GetValue(),
	}
	return out
}

func convertCommandResponse(c *rxlib.CommandResponse) *protoruntime.CommandResponse {
	out := &protoruntime.CommandResponse{
		SenderID:   c.SenderID,
		Count:      int32(nils.GetInt(c.Count)),
		MapStrings: c.MapStrings,
		Number:     nils.GetFloat64(c.Float),
		Boolean:    nils.GetBool(c.Bool),
		Error:      c.Error,
		ReturnType: c.ReturnType,
		//Response:   c.CommandResponse,
	}
	return out
}

func ConvertCommandResponse(c *rxlib.CommandResponse) *protoruntime.CommandResponse {
	cmd := convertCommandResponse(c)
	if len(c.CommandResponse) > 0 {
		var out []*protoruntime.CommandResponse
		for _, response := range c.CommandResponse {
			out = append(out, convertCommandResponse(response))
		}
		cmd.Response = out
	}
	if len(c.SerializeObjects) > 0 {
		var out []*protoruntime.Object
		for _, response := range c.SerializeObjects {
			out = append(out, ObjectConfigToProto(response))
		}
		cmd.SerializeObjects = out
	}
	return cmd
}
