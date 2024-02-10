package rxlib

type ObjectStatus string

const (
	StatsProcessing     ObjectStatus = "processing"      // processing its logic
	StatsLoaded         ObjectStatus = "loaded"          // the node Start() methods has at least been called once
	StatsLoopProcessing ObjectStatus = "loop-processing" // this would be used for a long-running process like a driver polling
	StatsHalted         ObjectStatus = "halted"          // critical Err so Obj processing has been disabled
	StatsEnabled        ObjectStatus = "enabled"         // disabled by the user
	StatsDisabled       ObjectStatus = "disabled"        // disabled by the user
	StatsIdle           ObjectStatus = "idle"            // idle is waiting for a new message to process
)

type ObjectType string

const (
	Logic               ObjectType = "logic"         // logic: this would be something like a math node
	Driver              ObjectType = "driver"        //driver: this would be something like a modbus network
	RubixNetwork        ObjectType = "rubix-network" //rubix-network: is a network a ROS instances
	Service             ObjectType = "service"       // service: this would be something like user service
	ExtensionObjectType ObjectType = "extension"     // extension: this would be something component of a service to an Obj; eg we have the history service, and we will add the COV component
)

var AllObjectType = []ObjectType{Logic, Driver, Service, ExtensionObjectType, RubixNetwork}
