package runtime

import (
	"github.com/NubeIO/rxlib"
	runtimeServer "github.com/NubeIO/rxlib/protos/runtime/runtimeserver"
)

func ObjectConfigToProto(obj *rxlib.ObjectConfig) *runtimeServer.Object {
	return &runtimeServer.Object{
		Id: obj.ID,
	}
}

func ObjectConfigFromProto(protoObj *runtimeServer.Object) *rxlib.ObjectConfig {
	return &rxlib.ObjectConfig{
		ID: protoObj.Id,
	}
}

func ConvertProtoConnectionToStruct(protoConn *runtimeServer.Connection) *rxlib.Connection {
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

func ConvertStructConnectionToProto(conn *rxlib.Connection) *runtimeServer.Connection {
	protoConn := &runtimeServer.Connection{
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
