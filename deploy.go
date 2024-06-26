package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/libs/restc"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/mitchellh/mapstructure"
)

type Deploy struct {
	Deleted []string                `json:"deleted"`
	New     []*runtime.ObjectConfig `json:"new"`
	Updated []*runtime.ObjectConfig `json:"updated"`
}

type Builder struct {
	UUID         string                `json:"uuid"`
	ObjectID     string                `json:"objectID"`
	Name         string                `json:"name"`
	ParentUUID   string                `json:"parentUUID"`
	ObjectConfig *runtime.ObjectConfig `json:"objectConfig"`
	Error        string                `json:"error"`
}

func (b *Builder) ToObject() *runtime.ObjectConfig {
	return newBuilderObject(b)
}

func (b *Builder) ToObjects() []*runtime.ObjectConfig {
	return []*runtime.ObjectConfig{newBuilderObject(b)}
}

func (inst *RuntimeImpl) ObjectBuilder(body *Builder) *Builder {
	return body
}

func newBuilderObject(body *Builder) *runtime.ObjectConfig {
	if body == nil {
		return nil
	}

	meta := body.ObjectConfig.GetMeta()
	if meta == nil {
		objectUUID := body.UUID
		if objectUUID == "" {
			objectUUID = helpers.UUID()
		}
		meta = &runtime.Meta{
			ObjectUUID: objectUUID,
			ParentUUID: body.ParentUUID,
			Position: &runtime.Position{
				PositionY: 0,
				PositionX: 0,
			},
		}
	}

	objectConfig := &runtime.ObjectConfig{
		Id:          body.ObjectID,
		Info:        body.ObjectConfig.GetInfo(),
		Inputs:      body.ObjectConfig.GetInputs(),
		Outputs:     body.ObjectConfig.GetOutputs(),
		Connections: body.ObjectConfig.GetConnections(),
		Settings:    body.ObjectConfig.GetSettings(),
		Stats:       body.ObjectConfig.GetStats(),
		Meta:        meta,
	}
	return objectConfig
}

func (inst *RuntimeImpl) UUID() string {
	return helpers.UUID()
}

type DeployResponse struct {
	Message string `json:"message"`
}

func (inst *RuntimeImpl) Deploy(body *Deploy) *DeployResponse {
	var invalidBody bool
	var message string
	if body == nil {
		invalidBody = true
		message = "body is nil"
	}
	if !invalidBody {
		if body.Deleted == nil && body.New == nil && body.Updated == nil {
			invalidBody = true
			message = "nothing to deploy"
		}
	}

	if invalidBody {
		var message = fmt.Sprintf("Deploy failed. %s", message)
		return &DeployResponse{
			Message: message,
		}
	}

	var existingCount = len(inst.Get())
	opts := &restc.Options{
		Headers: nil,
		Body:    body,
	}

	resp := inst.rest.Execute("POST", "http://localhost:1770/api/runtime", opts)
	var ok bool
	if resp.Code() >= 200 && resp.Code() < 300 {
		ok = true
	}
	if !ok {
		var message = fmt.Sprintf("Deploy failed. Response code: %d", resp.Code())
		if resp.GetError() != "" {
			message = fmt.Sprintf("Deploy failed. Response err: %s", resp.GetError())
		}
		return &DeployResponse{
			Message: message,
		}
	}

	var newCount = len(inst.Get())
	message = fmt.Sprintf("existingCount: %d current objects count: %d", existingCount, newCount)
	return &DeployResponse{
		Message: message,
	}
}

func (inst *RuntimeImpl) SerializeObjects(includePortValues bool, objects []Object) []*runtime.ObjectConfig {
	var serializedObjects []*runtime.ObjectConfig
	for _, object := range objects {
		serializedObjects = append(serializedObjects, inst.serializeObject(includePortValues, object))
	}
	return serializedObjects
}

func (inst *RuntimeImpl) serializeObject(includePortValues bool, object Object) *runtime.ObjectConfig {
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
	if includePortValues {
		objectConfig.PortValues = inst.GetObjectValues(object.GetUUID())
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
	p := &runtime.Port{
		Id:              obj.ID,
		Name:            obj.Name,
		PortUUID:        obj.UUID,
		Direction:       string(obj.Direction),
		DataType:        string(obj.DataType),
		DefaultPosition: int32(obj.DefaultPosition),
	}
	if obj.Transformation != nil {
		protoStruct, err := priority.ToProtoStruct(obj.Transformation)
		if err != nil {
			return p
		}
		p.Transformation = protoStruct
	}
	return p
}

func ProtoToPort(port *runtime.Port) *Port {
	p := &Port{
		ID:              port.Id,
		Name:            port.Name,
		UUID:            port.PortUUID,
		Transformation:  nil,
		Direction:       PortDirection(port.Direction),
		DataType:        priority.Type(port.DataType),
		DefaultPosition: int(port.DefaultPosition),
	}
	if port.GetTransformation() != nil {
		t, _ := convertTransformation(port)
		p.Transformation = t
	}
	return p
}

func ConvertTransformation(transformation *runtime.ObjectTransformations) (*priority.Transformations, error) {
	var trans *priority.Transformations
	err := mapstructure.Decode(transformation.GetTransformation().AsMap(), &trans)
	return trans, err
}

func convertTransformation(port *runtime.Port) (*priority.Transformations, error) {
	var trans *priority.Transformations
	err := mapstructure.Decode(port.GetTransformation().AsMap(), &trans)
	return trans, err
}
