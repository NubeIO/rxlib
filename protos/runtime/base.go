package runtime

import (
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/protos/runtime/protoruntime"
)

func ObjectConfigToProto(obj *rxlib.ObjectConfig) *protoruntime.Object {
	return &protoruntime.Object{
		Id:          obj.ID,
		Info:        ObjectInfoToProto(obj.Info),
		Inputs:      nil,
		Outputs:     nil,
		Meta:        nil,
		Stats:       nil,
		Connections: nil,
	}
}

func ObjectInfoToProto(info *rxlib.Info) *protoruntime.Info {
	if info == nil {
		return nil
	}

	protoObj := &protoruntime.Info{
		ObjectId:                 info.ObjectID,
		ObjectType:               string(info.ObjectType),
		Category:                 info.Category,
		PluginName:               info.PluginName,
		WorkingGroup:             info.WorkingGroup,
		WorkingGroupLeader:       info.WorkingGroupLeader,
		WorkingGroupParent:       info.WorkingGroupParent,
		WorkingGroupObjects:      info.WorkingGroupObjects,
		WorkingGroupChildObjects: info.WorkingGroupChildObjects,
		ObjectTags:               info.ObjectTags,
		Permissions: &protoruntime.Permissions{
			AllPermissions: info.Permissions.AllPermissions,
			CanBeCreated:   info.Permissions.CanBeCreated,
			CanBeDeleted:   info.Permissions.CanBeDeleted,
			CanBeUpdated:   info.Permissions.CanBeUpdated,
			ReadOnly:       info.Permissions.ReadOnly,
			AllowHotFix:    info.Permissions.AllowHotFix,
			ForceDelete:    info.Permissions.ForceDelete,
		},
		Requirements: &protoruntime.Requirements{
			CallResetOnDeploy:    info.Requirements.CallResetOnDeploy,
			AllowRuntimeAccess:   info.Requirements.AllowRuntimeAccess,
			MaxOne:               info.Requirements.MaxOne,
			MustLiveInObjectType: info.Requirements.MustLiveInObjectType,
			MustLiveParent:       info.Requirements.MustLiveParent,
			RequiresLogger:       info.Requirements.RequiresLogger,
			SupportsActions:      info.Requirements.SupportsActions,
			ServicesRequirements: info.Requirements.ServicesRequirements,
		},
	}

	return protoObj
}

func ObjectInfoFromProto(protoObj *protoruntime.Info) *rxlib.Info {
	if protoObj == nil {
		return nil
	}

	info := &rxlib.Info{
		ObjectID:                 protoObj.ObjectId,
		ObjectType:               rxlib.ObjectType(protoObj.ObjectType),
		Category:                 protoObj.Category,
		PluginName:               protoObj.PluginName,
		WorkingGroup:             protoObj.WorkingGroup,
		WorkingGroupLeader:       protoObj.WorkingGroupLeader,
		WorkingGroupParent:       protoObj.WorkingGroupParent,
		WorkingGroupObjects:      protoObj.WorkingGroupObjects,
		WorkingGroupChildObjects: protoObj.WorkingGroupChildObjects,
		ObjectTags:               protoObj.GetObjectTags(),
		Permissions: &rxlib.Permissions{
			AllPermissions: protoObj.Permissions.AllPermissions,
			CanBeCreated:   protoObj.Permissions.CanBeCreated,
			CanBeDeleted:   protoObj.Permissions.CanBeDeleted,
			CanBeUpdated:   protoObj.Permissions.CanBeUpdated,
			ReadOnly:       protoObj.Permissions.ReadOnly,
			AllowHotFix:    protoObj.Permissions.AllowHotFix,
			ForceDelete:    protoObj.Permissions.ForceDelete,
		},
		Requirements: &rxlib.Requirements{
			CallResetOnDeploy:    protoObj.Requirements.CallResetOnDeploy,
			AllowRuntimeAccess:   protoObj.Requirements.AllowRuntimeAccess,
			MaxOne:               protoObj.Requirements.MaxOne,
			MustLiveInObjectType: protoObj.Requirements.MustLiveInObjectType,
			MustLiveParent:       protoObj.Requirements.MustLiveParent,
			RequiresLogger:       protoObj.Requirements.RequiresLogger,
			SupportsActions:      protoObj.Requirements.SupportsActions,
			ServicesRequirements: protoObj.Requirements.ServicesRequirements,
			LoggerOpts:           &rxlib.LoggerOpts{
				// Fill in the fields accordingly
			},
		},
	}

	return info
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

func ConvertCommand(c *protoruntime.Command) *rxlib.Command {
	out := &rxlib.Command{
		TargetGlobalID:   c.GetTargetGlobalID(),
		SenderGlobalID:   c.GetSenderGlobalID(),
		SenderObjectUUID: c.GetSenderObjectUUID(),
		TransactionUUID:  c.GetTransactionUUID(),
		Key:              c.GetKey(),
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
		Any:        c.Any,
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
