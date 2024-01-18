package rxlib

type ObjectStatus string

const (
	StatsProcessing ObjectStatus = "processing" // processing its logic
	StatsHalted     ObjectStatus = "halted"     // critical error so object processing has been disabled
	StatsDisabled   ObjectStatus = "disabled"   // disabled by the user
	StatsIdle       ObjectStatus = "idle"       // idle is waiting for a new message to process
)

type ObjectType string

const (
	Logic            ObjectType = "logic"             // logic: this would be something like a math node
	Driver           ObjectType = "driver"            //driver: this would be something like a modbus network
	Service          ObjectType = "service"           // service: this would be something like user service
	Component        ObjectType = "component"         // component: this would be something component of a service to a object; eg we have the history service, and we will add the COV component
	MandatoryService ObjectType = "mandatory-service" // mandatoryService: this would be something like time service, and can not be added or removed
)

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

type ObjectRequirement string

const (
	AllowRuntimeAccess      ObjectRequirement = "allow-runtime-access" // for nodes that want/need access to all other nodes
	MaxOne                  ObjectRequirement = "max-one"              // maxOne: user can only add one. for example a plugin for setting IP/network address
	IsParent                ObjectRequirement = "is-parent"            // isParent: is a parent like a network
	IsChild                 ObjectRequirement = "is-child"             // isChild: is a child like a device
	HasChildren             ObjectRequirement = "has-children"         // hasChildren: like a modbus network has a device the device is the child
	SupportsAddingComponent ObjectRequirement = "supports-component"
	SupportsWebRoute        ObjectRequirement = "supports-router"
	SupportsDB              ObjectRequirement = "supports-db"
	DisableCreation         ObjectRequirement = "creation-disable"
	ReadOnly                ObjectRequirement = "read-only"
	ReadEdit                ObjectRequirement = "read-with-edit"
	ReadEditDelete          ObjectRequirement = "read-with-edit/delete"
)

func RequirementAllowRuntimeAccess() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        AllowRuntimeAccess,
		Description: "This will allow a node access to all the other nodes in the runtime",
	}
	return v
}

func RequirementDisableCreation() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        DisableCreation,
		Description: "user can not add a new object from the UI",
	}
	return v
}

// RequirementWebRouter is a read only ojbject meaning the UI can not edit it
func RequirementReadOnly() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        ReadOnly,
		Description: "object can not be deleted, edited",
	}
	return v
}

func RequirementReadEdit() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        ReadEdit,
		Description: "user can read and edit from the UI",
	}
	return v
}

func RequirementReadEditDelete() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        ReadEditDelete,
		Description: "user can read, edit and delete from the UI",
	}
	return v
}

// RequirementWebRouter is a read only ojbject meaning the UI can not edit it
func RequirementSupportsDB() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        SupportsDB,
		Description: "has databse acess on rx-server",
	}
	return v
}

// RequirementWebRouter needs a web router
func RequirementWebRouter() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        SupportsWebRoute,
		Description: "Has rest-api endpoints",
	}
	return v
}

// RequirementMaxOne this object can only be added once
func RequirementMaxOne() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        MaxOne,
		Description: "user can only add one. for example a plugin for setting IP/network address",
	}
	return v
}

// RequirementIsParent this object can have child objects
func RequirementIsParent() ObjectTypeRequirements {
	v := ObjectTypeRequirements{
		Type:        IsParent,
		Description: "IsParent: is a parent object like a network and a device would be a child object",
	}
	return v
}

type ObjectTypeRequirements struct {
	Type        ObjectRequirement `json:"type"`
	Description string            `json:"description"`
}
