package rxlib

import (
	"fmt"
)

type ObjectInfo interface {
	InfoBuilder
}

type permissions struct {
	AllPermissions bool `json:"allPermissions"`
	CanBeCreated   bool `json:"canBeCreated"`
	CanBeDeleted   bool `json:"canBeDeleted"`
	CanBeUpdated   bool `json:"canBeUpdated"`
	ReadOnly       bool `json:"readOnly"`
}

type requirements struct {
	SupportsWebRoute        bool `json:"supportsWebRoute"`
	AllowRuntimeAccess      bool `json:"allowRuntimeAccess"`
	MaxOne                  bool `json:"maxOne"`
	IsParent                bool `json:"isParent"`
	HasChildren             bool `json:"hasChildren"`
	SupportsAddingComponent bool `json:"supportsAddingComponent"`
}

type Info struct {
	ObjectID     string          `json:"objectID"`
	ObjectUUID   string          `json:"objectUUID"`
	ObjectName   string          `json:"objectName"`
	Category     string          `json:"category"`
	PluginName   string          `json:"pluginName"`
	Permissions  *permissions    `json:"permissions"`
	Requirements *requirements   `json:"requirements"`
	ObjectTags   []ObjectTypeTag `json:"objectTags"`
}

type InfoBuilder interface {
	Build() *Info
	String() string

	// id
	SetObjectID(objectID string) InfoBuilder

	// uuid
	SetObjectUUID(objectUUID string) InfoBuilder
	GetObjectUUID() string

	// name
	SetObjectName(objectName string) InfoBuilder
	GetObjectName() string

	// plugin
	SetPluginName(pluginName string) InfoBuilder
	GetPluginName() string

	// permissions
	GetPermissions() *permissions
	SetReadOnly() InfoBuilder
	SetAllPermissions() InfoBuilder
	SetCanBeCreated() InfoBuilder
	SetCanBeDeleted() InfoBuilder
	SetCanBeUpdated() InfoBuilder

	// requirements
	GetRequirements() *requirements
	SetSupportsWebRoute() InfoBuilder
	SetAllowRuntimeAccess() InfoBuilder
	SetMaxOne() InfoBuilder
	SetIsParent() InfoBuilder
	SetHasChildren() InfoBuilder
	SetSupportsAddingComponent() InfoBuilder

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

func (builder *infoBuilder) SetObjectID(objectID string) InfoBuilder {
	builder.info.ObjectID = objectID
	return builder
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

func (builder *infoBuilder) SetIsParent() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.IsParent = true
	return builder
}

func (builder *infoBuilder) SetHasChildren() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.HasChildren = true
	return builder
}

func (builder *infoBuilder) SetSupportsAddingComponent() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.SupportsAddingComponent = true
	return builder
}

func (builder *infoBuilder) SetReadOnly() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.ReadOnly = true
	builder.info.Permissions.AllPermissions = false
	builder.info.Permissions.CanBeCreated = false
	builder.info.Permissions.CanBeDeleted = false
	builder.info.Permissions.CanBeUpdated = false
	return builder
}

func (builder *infoBuilder) SetAllPermissions() InfoBuilder {
	ensurePermissions(builder.info)
	builder.info.Permissions.AllPermissions = true
	builder.info.Permissions.CanBeCreated = true
	builder.info.Permissions.CanBeDeleted = true
	builder.info.Permissions.CanBeUpdated = true
	builder.info.Permissions.ReadOnly = false
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

func (builder *infoBuilder) GetPermissions() *permissions {
	return builder.info.Permissions
}

func (builder *infoBuilder) GetRequirements() *requirements {
	return builder.info.Requirements
}

// SetObjectUUID sets the ObjectUUID for the InfoBuilder.
func (builder *infoBuilder) SetObjectUUID(objectUUID string) InfoBuilder {
	builder.info.ObjectUUID = objectUUID
	return builder
}

// SetObjectName sets the ObjectName for the InfoBuilder.
func (builder *infoBuilder) SetObjectName(objectName string) InfoBuilder {
	builder.info.ObjectName = objectName
	return builder
}

// GetObjectUUID returns the ObjectUUID associated with the InfoBuilder.
func (builder *infoBuilder) GetObjectUUID() string {
	return builder.info.ObjectUUID
}

// GetObjectName returns the ObjectName associated with the InfoBuilder.
func (builder *infoBuilder) GetObjectName() string {
	return builder.info.ObjectName
}

func (builder *infoBuilder) Build() *Info {
	return builder.info
}

func ensurePermissions(info *Info) {
	if info.Permissions == nil {
		info.Permissions = &permissions{}
	}
}

func ensureRequirements(info *Info) {
	if info.Requirements == nil {
		info.Requirements = &requirements{}
	}
}

func (builder *infoBuilder) String() string {
	return builder.info.String()
}

func (info *Info) String() string {
	return fmt.Sprintf("ObjectID: %s\nPluginName: %s\nCategory: %s\nPermissions: %+v", info.ObjectID, info.PluginName, info.Category, info.Permissions)
}
