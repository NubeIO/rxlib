package rxlib

import (
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

// Connection defines a structure for input subscriptions.
//type Connection struct {
//	UUID                 string        `json:"uuid,omitempty"` //the uuid of the rubix ***not needed for UI***
//	TargetConnectionUUID string        `json:"targetConnectionUUID,omitempty"`
//	SourceUUID           string        `json:"source,omitempty"` // will always be the output Obj
//	SourcePort           string        `json:"sourceHandle,omitempty"`
//	SourcePortUUID       string        `json:"sourcePortUUID,omitempty"` // output portID
//	TargetUUID           string        `json:"target,omitempty"`         // objectUUID that has the input rubix
//	TargetPort           string        `json:"targetHandle,omitempty"`   // input portID
//	TargetPortUUID       string        `json:"targetPortUUID,omitempty"`
//	IsExistingConnection bool          `json:"IsExistingConnection,omitempty"`
//	FlowDirection        FlowDirection `json:"flowDirection,omitempty"` // subscriber is if it's in an input and publisher or an output ***not needed for UI***
//	Disable              bool          `json:"disable,omitempty"`
//	IsError              bool          `json:"isError,omitempty"`
//	Created              *time.Time    `json:"created,omitempty"`
//	LastOk               *time.Time    `json:"LastOk,omitempty"`
//	LastFail             *time.Time    `json:"LastFail,omitempty"`
//	FailCount            int           `json:"failCount,omitempty"`
//	Error                []string      `json:"Err,omitempty"`
//}
//
//func (c *Connection) GetUUID() string {
//	return c.UUID
//}
//
//func (c *Connection) GetSourceUUID() string {
//	return c.SourceUUID
//}
//
//func (c *Connection) GetSourcePort() string {
//	return c.SourcePort
//}
//
//func (c *Connection) GetSourcePortUUID() string {
//	return c.SourcePortUUID
//}
//
//func (c *Connection) GetTargetUUID() string {
//	return c.TargetUUID
//}
//
//func (c *Connection) GetTargetPort() string {
//	return c.TargetPort
//}
//
//func (c *Connection) GetTargetPortUUID() string {
//	return c.TargetPortUUID
//}
//
//func (c *Connection) GetFlowDirection() FlowDirection {
//	return c.FlowDirection
//}
//
//func (c *Connection) DirectionPublisher() bool {
//	if c.FlowDirection == DirectionPublisher {
//		return true
//	}
//	return false
//}
//
//func (c *Connection) DirectionSubscriber() bool {
//	if c.FlowDirection == DirectionSubscriber {
//		return true
//	}
//	return false
//}

/*
Example of a Trigger Obj output connected to a Count Obj input
This is what's needed for the UI to work

 Trigger Obj (output)
 "source": "triggerABC",
 "sourceHandle": "output",
 "target": "counterABC",
 "targetHandle": "input",
 "flowDirection": "publisher"

 Count Obj (input)
"source": "triggerABC",
"sourceHandle": "output",
"target": "counterABC",
"targetHandle": "input",
"flowDirection": "subscriber"

*/

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
