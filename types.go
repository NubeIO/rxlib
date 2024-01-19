package rxlib

type ObjectStatus string

const (
	StatsProcessing     ObjectStatus = "processing"      // processing its logic
	StatsLoaded         ObjectStatus = "loaded"          // the node Start() methods has at least been called once
	StatsLoopProcessing ObjectStatus = "loop-processing" // this would be used for a long running process like a driver polling
	StatsHalted         ObjectStatus = "halted"          // critical error so object processing has been disabled
	StatsEnabled        ObjectStatus = "enabled"         // disabled by the user
	StatsDisabled       ObjectStatus = "disabled"        // disabled by the user
	StatsIdle           ObjectStatus = "idle"            // idle is waiting for a new message to process
)

type ObjectType string

const (
	Logic     ObjectType = "logic"     // logic: this would be something like a math node
	Driver    ObjectType = "driver"    //driver: this would be something like a modbus network
	Service   ObjectType = "service"   // service: this would be something like user service
	Component ObjectType = "component" // component: this would be something component of a service to a object; eg we have the history service, and we will add the COV component
)

var AllObjectType = []ObjectType{Logic, Driver, Service, Component}

type ObjectTypeTag string

const (
	DriversTag    ObjectTypeTag = "drivers"
	BACnetTag     ObjectTypeTag = "bacnet"
	ModbusTag     ObjectTypeTag = "modbus"
	LogicTag      ObjectTypeTag = "logic"
	TimeTag       ObjectTypeTag = "time"
	DateTag       ObjectTypeTag = "date"
	ServicesTag   ObjectTypeTag = "service"
	UsersTag      ObjectTypeTag = "users"
	NetworkingTag ObjectTypeTag = "networking"
	IpAddressTag  ObjectTypeTag = "ip-address"
)
