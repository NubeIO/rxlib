package rxlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"log"
)

type ObjectInfo interface {
	InfoBuilder
}

type RubixRequirement string

const (
	RubixGinRouter        RubixRequirement = "web-router"
	RubixAlarmsManager    RubixRequirement = "alarms-manager"
	RubixSchedulesManager RubixRequirement = "schedules-manager"
)

type InfoBuilder interface {
	Build() *runtime.Info
	String() string

	// id
	SetID(objectID string) InfoBuilder
	GetID() string

	// Obj type is for example a driver, service, logic
	SetObjectType(objectType ObjectType) InfoBuilder
	GetObjectType() ObjectType

	// category
	SetCategory(value string) InfoBuilder
	GetCategory() string

	// working group is a grough of Obj that internally work together; for example we have MQTT broker and sub Obj would be a working group
	SetWorkingGroup(value string) InfoBuilder
	GetWorkingGroup() string

	SetWorkingGroupLeader(value string) InfoBuilder
	GetWorkingGroupLeader() string

	SetWorkingGroupObjects(value ...string) InfoBuilder
	GetWorkingGroupObjects() []string
	SetWorkingGroupChildObjects(value ...string) InfoBuilder
	GetWorkingGroupChildObjects() []string
	SetWorkingGroupParent(value string) InfoBuilder
	GetWorkingGroupParent() string

	// plugin
	SetPluginName(pluginName string) InfoBuilder
	GetPluginName() string

	// permissions
	GetPermissions() *runtime.Permissions
	SetReadOnly() InfoBuilder
	SetAllPermissions() InfoBuilder
	SetCanBeCreated() InfoBuilder
	SetCanBeDeleted() InfoBuilder
	SetCanBeUpdated() InfoBuilder

	// requirements
	GetRequirements() *runtime.Requirements
	SetCallResetOnDeploy() InfoBuilder
	SetAllowRuntimeAccess() InfoBuilder
	SetMaxOne() InfoBuilder // only max one Obj can be added
	SetLogger(opts []string) InfoBuilder

	SetSupportsActions() InfoBuilder
	GetSupportsActions() bool

	SetMustLiveInObjectType() InfoBuilder
	GetMustLiveInObjectType() bool
	SetMustLiveParent() InfoBuilder
	GetMustLiveParent() bool

	// tags
	AddObjectTags(objectTypeTag ...string) InfoBuilder
	GetObjectTags() []string
}

func NewObjectInfo() InfoBuilder {
	return &infoBuilder{info: &runtime.Info{}}
}

type infoBuilder struct {
	info *runtime.Info
}

func (builder *infoBuilder) Build() *runtime.Info {
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
	builder.info.ObjectType = string(objectType)
	return builder
}

func (builder *infoBuilder) GetObjectType() ObjectType {
	return ObjectType(builder.info.ObjectType)
}

func (builder *infoBuilder) SetCategory(value string) InfoBuilder {
	builder.info.Category = value
	return builder
}

func (builder *infoBuilder) GetCategory() string {
	return builder.info.Category
}

func (builder *infoBuilder) SetWorkingGroup(value string) InfoBuilder {
	builder.info.WorkingGroup = value
	return builder
}

func (builder *infoBuilder) GetWorkingGroup() string {
	return builder.info.WorkingGroup
}

func (builder *infoBuilder) SetWorkingGroupLeader(value string) InfoBuilder {
	builder.info.WorkingGroupLeader = value
	return builder
}

func (builder *infoBuilder) GetWorkingGroupLeader() string {
	return builder.info.WorkingGroupLeader
}

func (builder *infoBuilder) SetWorkingGroupParent(value string) InfoBuilder {
	builder.info.WorkingGroupParent = value
	return builder
}

func (builder *infoBuilder) GetWorkingGroupParent() string {
	return builder.info.WorkingGroupParent
}

func (builder *infoBuilder) SetWorkingGroupObjects(value ...string) InfoBuilder {
	builder.info.WorkingGroupObjects = append(builder.info.WorkingGroupObjects, value...)
	return builder
}

func (builder *infoBuilder) GetWorkingGroupObjects() []string {
	return builder.info.WorkingGroupObjects
}

func (builder *infoBuilder) SetWorkingGroupChildObjects(value ...string) InfoBuilder {
	builder.info.WorkingGroupChildObjects = append(builder.info.WorkingGroupChildObjects, value...)
	return builder
}

func (builder *infoBuilder) GetWorkingGroupChildObjects() []string {
	return builder.info.WorkingGroupChildObjects
}

func (builder *infoBuilder) SetLogger(opts []string) InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.RequiresLogger = true
	if opts == nil {
		log.Fatalf("rxlib.SetLogger opts can not be empty")
	}
	return builder
}

func (builder *infoBuilder) SetCallResetOnDeploy() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.CallResetOnDeploy = true
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

func (builder *infoBuilder) SetSupportsActions() InfoBuilder {
	ensureRequirements(builder.info)
	builder.info.Requirements.SupportsActions = true
	return builder
}

func (builder *infoBuilder) GetSupportsActions() bool {
	return builder.info.Requirements.SupportsActions
}

func (builder *infoBuilder) SetPluginName(pluginName string) InfoBuilder {
	builder.info.PluginName = pluginName
	return builder
}

func (builder *infoBuilder) GetPluginName() string {
	return builder.info.PluginName
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

func (builder *infoBuilder) AddObjectTags(objectTypeTag ...string) InfoBuilder {
	builder.info.ObjectTags = append(builder.info.ObjectTags, objectTypeTag...)
	return builder
}

// GetObjectTags returns the ObjectTypeTags associated with the InfoBuilder.
func (builder *infoBuilder) GetObjectTags() []string {
	return builder.info.ObjectTags
}

func (builder *infoBuilder) GetPermissions() *runtime.Permissions {
	return builder.info.Permissions
}

func (builder *infoBuilder) GetRequirements() *runtime.Requirements {
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

}

func ensurePermissions(info *runtime.Info) {
	if info.Permissions == nil {
		info.Permissions = &runtime.Permissions{}
	}
}

func ensureRequirements(info *runtime.Info) {
	if info.Requirements == nil {
		info.Requirements = &runtime.Requirements{}
	}
}

func (builder *infoBuilder) String() string {
	return builder.info.String()
}

//func (info *Info) String() string {
//	return fmt.Sprintf("ObjectID: %s\nPluginName: %s\nCategory: %s\nPermissions: %+v", info.ObjectID, info.PluginName, info.Category, info.Permissions)
//}

func crashMe(name string) {
	log.Fatalf("rxlib.Checks() %s is empty", name)
}
