package runtime

import (
	"github.com/NubeIO/rxlib"
	protoruntime "github.com/NubeIO/rxlib/protos/runtime/pb"
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
