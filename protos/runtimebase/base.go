package runtimebase

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
