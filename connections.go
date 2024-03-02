package rxlib

import (
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

func NewConnection(sourceUUID, sourcePort, targetUUID, targetPort string) (publisher *runtime.Connection, subscriber *runtime.Connection) {
	sourceConnectionUUID := helpers.UUID()
	targetConnectionUUID := helpers.UUID()
	publisher = &runtime.Connection{
		ConnectionUUID:       sourceConnectionUUID,
		TargetConnectionUUID: targetConnectionUUID,
		SourceUUID:           sourceUUID,
		SourcePort:           sourcePort,
		TargetUUID:           targetUUID,
		TargetPort:           targetPort,
		FlowDirection:        DirectionPublisher,
	}
	subscriber = &runtime.Connection{
		ConnectionUUID: targetConnectionUUID,
		SourceUUID:     sourceUUID,
		SourcePort:     sourcePort,
		TargetUUID:     targetUUID,
		TargetPort:     targetPort,
		FlowDirection:  DirectionSubscriber,
	}
	return publisher, subscriber
}

//type UpdateConnectionsReport struct {
//	ExistingCount int      `json:"existingCount"` // before we started updating/deleting get the existing count
//	DeletedCount  int      `json:"deletedCount"`
//	DeployedCount int      `json:"newCount"`
//	Errors        []string `json:"errors"`
//}

//type RemoveConnectionReport struct {
//	ConnectionUUID string        `json:"connectionUUID,omitempty"`
//	ObjectUUID     string        `json:"objectUUID,omitempty"`
//	ObjectID       string        `json:"objectID,omitempty"`
//	TargetUUID     string        `json:"targetUUID,omitempty"`
//	TargetPort     string        `json:"targetPort,omitempty"`
//	SourceUUID     string        `json:"sourceUUID,omitempty"`
//	SourcePort     string        `json:"sourcePort,omitempty"`
//	FlowDirection  FlowDirection `json:"flowDirection,omitempty"`
//	Error          string        `json:"Err,omitempty"`
//}
