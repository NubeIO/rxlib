package rxlib

import (
	"fmt"
	"log"
)

type ObjectInfo interface {
	InfoBuilder
}

type permissions struct {
	AllPermissions bool `json:"allPermissions,omitempty"`
	CanBeCreated   bool `json:"canBeCreated,omitempty"`
	CanBeDeleted   bool `json:"canBeDeleted,omitempty"`
	CanBeUpdated   bool `json:"canBeUpdated,omitempty"`
	ReadOnly       bool `json:"readOnly,omitempty"`
}

type requirements struct {
	SupportsWebRoute        bool `json:"supportsWebRoute,omitempty"`
	AllowRuntimeAccess      bool `json:"allowRuntimeAccess,omitempty"`
	MaxOne                  bool `json:"maxOne,omitempty"`
	IsParent                bool `json:"isParent,omitempty"`
	HasChildren             bool `json:"hasChildren,omitempty"`
	SupportsAddingComponent bool `json:"supportsAddingComponent,omitempty"`
}

type Info struct {
	Settings     *Settings       `json:"settings"`
	ObjectID     string          `json:"objectID"`
	ObjectUUID   string          `json:"objectUUID"`
	ObjectName   string          `json:"objectName"`
	ParentUUID   string          `json:"parentUUID,omitempty"`
	Category     string          `json:"category"`
	PluginName   string          `json:"pluginName"`
	Permissions  *permissions    `json:"permissions"`
	Requirements *requirements   `json:"requirements"`
	ObjectTags   []ObjectTypeTag `json:"objectTags,omitempty"`
	Meta         *Meta           `json:"meta"`
}

type InfoBuilder interface {
	Build() *Info
	String() string

	// id
	SetID(objectID string) InfoBuilder
	GetID() string

	// uuid
	GetUUID() string

	// name
	GetName() string

	// category
	GetCategory() string

	// plugin
	SetPluginName(pluginName string) InfoBuilder
	GetPluginName() string

	// settings
	SetSettings(settings *Settings) InfoBuilder
	GetSettings() *Settings

	// meta, meta will also set the object-name at parentUUID
	SetMeta(meta *Meta) InfoBuilder
	GetMeta() *Meta

	// parent uuid
	GetParentUUID() string

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

func (builder *infoBuilder) SetID(objectID string) InfoBuilder {
	builder.info.ObjectID = objectID
	return builder
}
func (builder *infoBuilder) GetID() string {
	return builder.info.ObjectID
}

func (builder *infoBuilder) GetUUID() string {
	return builder.info.ObjectUUID
}

func (builder *infoBuilder) GetName() string {
	return builder.info.ObjectName
}

func (builder *infoBuilder) GetCategory() string {
	return builder.info.Category
}

func (builder *infoBuilder) SetPluginName(pluginName string) InfoBuilder {
	builder.info.PluginName = pluginName
	return builder
}

func (builder *infoBuilder) SetSettings(settings *Settings) InfoBuilder {
	builder.info.Settings = settings
	return builder
}

func (builder *infoBuilder) GetSettings() *Settings {
	return builder.info.Settings
}

func (builder *infoBuilder) SetMeta(meta *Meta) InfoBuilder {
	if meta == nil {
		log.Fatal("rxlib.SetMeta() meta is nil")
	}
	builder.info.Meta = meta
	builder.info.ObjectName = meta.ObjectName
	builder.info.ParentUUID = meta.ParentUUID
	return builder
}

func (builder *infoBuilder) GetMeta() *Meta {
	return builder.info.Meta
}

func (builder *infoBuilder) GetParentUUID() string {
	return builder.info.ParentUUID
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

func (builder *infoBuilder) checks() {
	// checks
	if builder.info.PluginName == "" {
		crashMe("info.PluginName")
	}
	if builder.info.ObjectUUID == "" {
		crashMe("info.ObjectUUID")
	}
	if builder.info.ObjectID == "" {
		crashMe("info.ObjectID")
	}
	if builder.info.Category == "" {
		crashMe("info.Category")
	}
}

func (builder *infoBuilder) Build() *Info {
	builder.checks()
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

func crashMe(name string) {
	log.Fatalf("rxlib.Checks() %s is empty", name)
}
