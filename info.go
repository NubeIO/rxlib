package rxlib

import (
	"fmt"
	"log"
)

type ObjectInfo interface {
	InfoBuilder
}

type Permissions struct {
	AllPermissions bool `json:"allPermissions,omitempty"`
	CanBeCreated   bool `json:"canBeCreated,omitempty"`
	CanBeDeleted   bool `json:"canBeDeleted,omitempty"`
	CanBeUpdated   bool `json:"canBeUpdated,omitempty"`
	ReadOnly       bool `json:"readOnly,omitempty"`
	AllowHotFix    bool `json:"allowHotFix,omitempty"`
	ForceDelete    bool `json:"forceDelete,omitempty"`
}

type Requirements struct {
	SupportsWebRoute     bool `json:"supportsWebRoute,omitempty"`
	AllowRuntimeAccess   bool `json:"allowRuntimeAccess,omitempty"`
	MaxOne               bool `json:"maxOne,omitempty"`
	HasChildObjects      bool `json:"hasChildObjects,omitempty"` // math object that has none, but say a modbus network will have childs; eg drives/points
	MustLiveInObjectType bool `json:"mustLiveInObjectType"`      // modbus-network can only be in object-type: drivers
	MustLiveParent       bool `json:"mustLiveParent"`            // a modbus device can only be added under its parent being a modbus-network
	RequiresLogger       bool `json:"requiresLogger"`
}

type Info struct {
	ObjectID     string          `json:"id"`
	ObjectType   ObjectType      `json:"type"`
	Category     string          `json:"category"`
	PluginName   string          `json:"pluginName"`
	Permissions  *Permissions    `json:"permissions"`
	Requirements *Requirements   `json:"requirements,omitempty"`
	ObjectTags   []ObjectTypeTag `json:"objectTags,omitempty"`
}

type InfoBuilder interface {
	Build() *Info
	String() string

	// id
	SetID(objectID string) InfoBuilder
	GetID() string

	// object type is for example a driver, service, logic
	SetObjectType(objectType ObjectType) InfoBuilder
	GetObjectType() ObjectType

	// category
	SetCategory(value string) InfoBuilder
	GetCategory() string

	// plugin
	SetPluginName(pluginName string) InfoBuilder
	GetPluginName() string

	// permissions
	GetPermissions() *Permissions
	SetReadOnly() InfoBuilder
	SetAllPermissions() InfoBuilder
	SetCanBeCreated() InfoBuilder
	SetCanBeDeleted() InfoBuilder
	SetCanBeUpdated() InfoBuilder

	// requirements
	GetRequirements() *Requirements
	SetSupportsWebRoute() InfoBuilder
	SetAllowRuntimeAccess() InfoBuilder
	SetMaxOne() InfoBuilder
	SetHasChildObjects() InfoBuilder
	SetLogger() InfoBuilder

	SetMustLiveInObjectType() InfoBuilder
	GetMustLiveInObjectType() bool
	SetMustLiveParent() InfoBuilder
	GetMustLiveParent() bool

	// tags
	AddObjectTags(objectTypeTag ...ObjectTypeTag)
	GetObjectTags() []ObjectTypeTag
}

func NewObjectInfo() InfoBuilder {
	return &infoBuilder{info: &Info{}}
}

type infoBuilder struct {
	info *Info
}

func (builder *infoBuilder) Build() *Info {
	builder.checks()
	return builder.info
}

func (builder *infoBuilder) SetID(objectID string) InfoBuilder {
	builder.info.ObjectID = objectID
	return builder
}

func (builder *infoBuilder) GetID() string {
	return builder.info.ObjectID
}

func (builder *infoBuilder) SetObjectType(objectType ObjectType) InfoBuilder {
	builder.info.ObjectType = objectType
	return builder
}

func (builder *infoBuilder) GetObjectType() ObjectType {
	return builder.info.ObjectType
}

func (builder *infoBuilder) SetCategory(value string) InfoBuilder {
	builder.info.Category = value
	return builder
}

func (builder *infoBuilder) SetLogger() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.RequiresLogger = true
	return builder
}

func (builder *infoBuilder) SetMustLiveInObjectType() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.MustLiveInObjectType = true
	return builder
}

func (builder *infoBuilder) GetMustLiveInObjectType() bool {
	return builder.info.Requirements.MustLiveInObjectType
}

func (builder *infoBuilder) SetMustLiveParent() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.MustLiveParent = true
	return builder
}

func (builder *infoBuilder) GetMustLiveParent() bool {
	return builder.info.Requirements.MustLiveParent
}

func (builder *infoBuilder) GetCategory() string {
	return builder.info.Category
}

func (builder *infoBuilder) SetPluginName(pluginName string) InfoBuilder {
	builder.info.PluginName = pluginName
	return builder
}

func (builder *infoBuilder) GetPluginName() string {
	return builder.info.PluginName
}

func (builder *infoBuilder) SetSupportsWebRoute() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.SupportsWebRoute = true
	return builder
}

func (builder *infoBuilder) SetAllowRuntimeAccess() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.AllowRuntimeAccess = true
	return builder
}

func (builder *infoBuilder) SetMaxOne() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.MaxOne = true
	return builder
}

func (builder *infoBuilder) SetHasChildObjects() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.HasChildObjects = true
	return builder
}

func (builder *infoBuilder) SetReadOnly() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.ReadOnly = true
	return builder
}

func (builder *infoBuilder) SetAllPermissions() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.AllPermissions = true
	return builder
}

func (builder *infoBuilder) SetCanBeCreated() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.CanBeCreated = true
	return builder
}

func (builder *infoBuilder) SetCanBeDeleted() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.CanBeDeleted = true
	return builder
}

func (builder *infoBuilder) SetCanBeUpdated() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.CanBeUpdated = true
	return builder
}

func (builder *infoBuilder) AddObjectTags(objectTypeTag ...ObjectTypeTag) {
	builder.info.ObjectTags = append(builder.info.ObjectTags, objectTypeTag...)
}

// GetObjectTags returns the ObjectTypeTags associated with the InfoBuilder.
func (builder *infoBuilder) GetObjectTags() []ObjectTypeTag {
	return builder.info.ObjectTags
}

func (builder *infoBuilder) GetPermissions() *Permissions {
	return builder.info.Permissions
}

func (builder *infoBuilder) GetRequirements() *Requirements {
	return builder.info.Requirements
}

func (builder *infoBuilder) checks() {
	// checks
	if builder.info.PluginName == "" {
		crashMe("info.PluginName")
	}
	if builder.info.ObjectID == "" {
		crashMe("info.ObjectID")
	}
	if builder.info.Category == "" {
		crashMe("info.Category")
	}
	if builder.info.ObjectType == "" {
		crashMe("info.ObjectType")
	}
	var validType bool
	for _, o := range AllObjectType {
		if o == builder.info.ObjectType {
			validType = true
		}
	}
	if !validType {
		log.Fatalf("rxlib.SetObjectType() invaild object type: %s try: %s", builder.info.ObjectType, AllObjectType[0])
	}

}

func ensurePermissions(info *Info) {
	if info.Permissions == nil {
		info.Permissions = &Permissions{}
	}
}

func ensureRequirements(info *Info) {
	if info.Requirements == nil {
		info.Requirements = &Requirements{}
	}
}

func (builder *infoBuilder) String() string {
	return builder.info.String()
}

func (info *Info) String() string {
	return fmt.Sprintf("ObjectID: %s\nPluginName: %s\nCategory: %s\nPermissions: %+v", info.ObjectID, info.PluginName, info.Category, info.Permissions)
}

func crashMe(name string) {
	log.Fatalf("rxlib.Checks() %s is empty", name)
}
