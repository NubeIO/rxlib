package rxlib

import (
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type MultipleConnection struct {
	InputPortID  string
	OutputPortID string
	IsOutput     bool
}

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
