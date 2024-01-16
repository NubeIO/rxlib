package rxlib

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
	Networking ObjectTypeTag = "networking"
	IpAddress  ObjectTypeTag = "ip-address"
)

type ObjectTypeRequirement string

const (
	MaxOne                  ObjectTypeRequirement = "max-one"      // maxOne: user can only add one. for example a plugin for setting IP/network address
	IsParent                ObjectTypeRequirement = "is-parent"    // isParent: is a parent like a network
	IsChild                 ObjectTypeRequirement = "is-child"     // isChild: is a child like a device
	HasChildren             ObjectTypeRequirement = "has-children" // hasChildren: like a modbus network has a device the device is the child
	SupportsAddingComponent ObjectTypeRequirement = "supports-component"
	SupportsWebRoute        ObjectTypeRequirement = "supports-router"
)

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
	Type        ObjectTypeRequirement `json:"type"`
	Description string                `json:"description"`
}
