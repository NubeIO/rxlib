package rxlib

// Connection defines a structure for input subscriptions.
type Connection struct {
	UUID                 string        `json:"uuid"`           //the uuid of the connection ***not needed for UI***
	SourceUUID           string        `json:"source"`         // will always be the output object
	SourcePort           string        `json:"sourceHandle"`   // output portID
	SourcePortUUID       string        `json:"sourcePortUUID"` // output portUUID ***not needed for UI***
	TargetUUID           string        `json:"target"`         // objectUUID that has the input connection
	TargetPort           string        `json:"targetHandle"`   // input portID
	TargetPortUUID       string        `json:"targetPortUUID"` // input portUUID ***not needed for UI***
	FlowDirection        FlowDirection `json:"flowDirection"`  // subscriber is if it's in an input and publisher or an output ***not needed for UI***
	IsExistingConnection bool          `json:"IsExistingConnection"`
}

/*
Example of a Trigger object output connected to a Count object input
This is what's needed for the UI to work

 Trigger object (output)
 "source": "triggerABC",
 "sourceHandle": "output",
 "target": "counterABC",
 "targetHandle": "input",
 "flowDirection": "publisher"

 Count object (input)
"source": "triggerABC",
"sourceHandle": "output",
"target": "counterABC",
"targetHandle": "input",
"flowDirection": "subscriber"

*/

func NewConnection(sourceUUID, sourcePort, targetUUID, targetPort string) (publisher *Connection, subscriber *Connection, err error) {
	publisher = &Connection{
		SourceUUID:    sourceUUID,
		SourcePort:    sourcePort,
		TargetUUID:    targetUUID,
		TargetPort:    targetPort,
		FlowDirection: DirectionPublisher,
	}
	subscriber = &Connection{
		SourceUUID:    publisher.SourceUUID,
		SourcePort:    publisher.SourcePort,
		TargetUUID:    publisher.TargetUUID,
		TargetPort:    publisher.TargetPort,
		FlowDirection: DirectionSubscriber,
	}
	return publisher, subscriber, nil
}

type UpdateConnectionsReport struct {
	ExistingCount int      `json:"existingCount"` // before we started updating/deleting get the existing count
	DeletedCount  int      `json:"deletedCount"`
	DeployedCount int      `json:"newCount"`
	Errors        []string `json:"errors"`
}

type RemoveConnectionReport struct {
	ConnectionUUID string        `json:"connectionUUID,omitempty"`
	ObjectUUID     string        `json:"objectUUID,omitempty"`
	ObjectID       string        `json:"objectID,omitempty"`
	TargetUUID     string        `json:"targetUUID,omitempty"`
	TargetPort     string        `json:"targetPort,omitempty"`
	SourceUUID     string        `json:"sourceUUID,omitempty"`
	SourcePort     string        `json:"sourcePort,omitempty"`
	FlowDirection  FlowDirection `json:"flowDirection,omitempty"`
	Error          string        `json:"error,omitempty"`
}
