package ndocs

var ObjectString = `
[
  {
    "name": "New",
    "description": "",
    "args": [
      "object Object",
      "opts ...any"
    ],
    "return": [
      "Object"
    ],
    "help": ""
  },
  {
    "name": "Init",
    "description": "",
    "args": null,
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "Start",
    "description": "",
    "args": null,
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "SetLoaded",
    "description": "used normally for the Start() to set it that it has booted",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "IsNotLoaded",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "IsLoaded",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "CommandList",
    "description": "",
    "args": null,
    "return": [
      "[]*Invoke"
    ],
    "help": ""
  },
  {
    "name": "Process",
    "description": "",
    "args": null,
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "Reset",
    "description": "",
    "args": null,
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "AllowsReset",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "Delete",
    "description": "",
    "args": null,
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "Lock",
    "description": "",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "Unlock",
    "description": "",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "IsLocked",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "IsUnlocked",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "InvokePayload",
    "description": "",
    "args": [
      "p *payload.Payload"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "Invoke",
    "description": "",
    "args": [
      "command *ExtendedCommand"
    ],
    "return": [
      "*CommandResponse",
      "error"
    ],
    "help": ""
  },
  {
    "name": "CommandObject",
    "description": "",
    "args": [
      "command *ExtendedCommand"
    ],
    "return": [
      "*CommandResponse"
    ],
    "help": ""
  },
  {
    "name": "Command",
    "description": "",
    "args": [
      "command *ExtendedCommand"
    ],
    "return": [
      "*runtime.CommandResponse"
    ],
    "help": ""
  },
  {
    "name": "CommandResponse",
    "description": "",
    "args": [
      "response *runtime.CommandResponse"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "AddRuntime",
    "description": "",
    "args": [
      "r Runtime"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "Runtime",
    "description": "",
    "args": null,
    "return": [
      "Runtime"
    ],
    "help": ""
  },
  {
    "name": "RemoveObjectFromRuntime",
    "description": "",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "GetParentObject",
    "description": "",
    "args": null,
    "return": [
      "Object"
    ],
    "help": ""
  },
  {
    "name": "GetParentUUID",
    "description": "",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "AddExtension",
    "description": "AddExtension extension are a way to extend the functionalists of an Obj; for example add a history extension",
    "args": [
      "extension Object"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetExtensions",
    "description": "",
    "args": null,
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "GetExtension",
    "description": "",
    "args": [
      "id string"
    ],
    "return": [
      "Object",
      "error"
    ],
    "help": ""
  },
  {
    "name": "DeleteExtension",
    "description": "",
    "args": [
      "name string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "SetRequiredExtensions",
    "description": "",
    "args": [
      "extension []*Extension"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetRequiredExtensions",
    "description": "",
    "args": null,
    "return": [
      "[]*Extension"
    ],
    "help": ""
  },
  {
    "name": "RequiredExtensionListCount",
    "description": "",
    "args": null,
    "return": [
      "int"
    ],
    "help": ""
  },
  {
    "name": "IsExtensionsAdded",
    "description": "",
    "args": [
      "objectID string"
    ],
    "return": [
      "int"
    ],
    "help": ""
  },
  {
    "name": "GetRequiredExtensionByName",
    "description": "",
    "args": [
      "extensionName string"
    ],
    "return": [
      "*Extension"
    ],
    "help": ""
  },
  {
    "name": "NewPort",
    "description": "ports",
    "args": [
      "port *Port"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "NewInputPort",
    "description": "",
    "args": [
      "port *NewPort"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "NewInputPorts",
    "description": "",
    "args": [
      "port []*NewPort"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "NewOutputPort",
    "description": "",
    "args": [
      "port *NewPort"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "NewOutputPorts",
    "description": "",
    "args": [
      "port []*NewPort"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetAllPorts",
    "description": "",
    "args": null,
    "return": [
      "[]*Port"
    ],
    "help": ""
  },
  {
    "name": "GetPortPayload",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "*payload.Payload",
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetPortValue",
    "description": "SetPortPayload(portID string, data *payload.Payload) error",
    "args": [
      "portID string"
    ],
    "return": [
      "*runtime.PortValue"
    ],
    "help": ""
  },
  {
    "name": "EnablePort",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "DisablePort",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "IsPortDisable",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "bool",
      "error"
    ],
    "help": ""
  },
  {
    "name": "AddAllTransformations",
    "description": "",
    "args": [
      "inputs []*Port",
      "outputs []*Port"
    ],
    "return": [
      "[]error"
    ],
    "help": ""
  },
  {
    "name": "OverrideValue",
    "description": "",
    "args": [
      "value any",
      "portID string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "ReleaseOverride",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "CreateConnection",
    "description": "",
    "args": [
      "connection *runtime.Connection"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "NewOutputConnection",
    "description": "",
    "args": [
      "portID string",
      "targetUUID string",
      "targetPort string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetPortConnections",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "[]*runtime.Connection"
    ],
    "help": ""
  },
  {
    "name": "GetConnection",
    "description": "",
    "args": [
      "uuid string"
    ],
    "return": [
      "*runtime.Connection"
    ],
    "help": ""
  },
  {
    "name": "GetExistingConnection",
    "description": "",
    "args": [
      "sourceObjectUUID string",
      "targetObjectUUID string",
      "targetPortID string"
    ],
    "return": [
      "*runtime.Connection"
    ],
    "help": ""
  },
  {
    "name": "GetConnections",
    "description": "",
    "args": null,
    "return": [
      "[]*runtime.Connection"
    ],
    "help": ""
  },
  {
    "name": "PortHasConnection",
    "description": "",
    "args": [
      "portID string"
    ],
    "return": [
      "bool",
      "int"
    ],
    "help": ""
  },
  {
    "name": "RemoveConnection",
    "description": "",
    "args": [
      "connection *runtime.Connection"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "DropConnections",
    "description": "",
    "args": null,
    "return": [
      "[]error"
    ],
    "help": ""
  },
  {
    "name": "RemoveOldConnections",
    "description": "",
    "args": [
      "newConnections []*runtime.Connection"
    ],
    "return": [
      "[]error"
    ],
    "help": ""
  },
  {
    "name": "AddSubscriptionConnection",
    "description": "",
    "args": [
      "sourceObjectUUID string",
      "sourcePortID string",
      "targetObjectUUID string",
      "targetPortID string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetInput",
    "description": "inputs",
    "args": [
      "id string"
    ],
    "return": [
      "*Port"
    ],
    "help": ""
  },
  {
    "name": "InputExists",
    "description": "",
    "args": [
      "id string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetInputs",
    "description": "",
    "args": null,
    "return": [
      "[]*Port"
    ],
    "help": ""
  },
  {
    "name": "GetInputByConnection",
    "description": "",
    "args": [
      "sourceObjectUUID string",
      "outputPortID string"
    ],
    "return": [
      "*Port"
    ],
    "help": ""
  },
  {
    "name": "GetInputByConnections",
    "description": "",
    "args": [
      "sourceObjectUUID string",
      "outputPortID string"
    ],
    "return": [
      "[]*Port"
    ],
    "help": ""
  },
  {
    "name": "UpdateInputsValue",
    "description": "",
    "args": [
      "payload *payload.Payload"
    ],
    "return": [
      "[]error"
    ],
    "help": ""
  },
  {
    "name": "GetOutputs",
    "description": "",
    "args": null,
    "return": [
      "[]*Port"
    ],
    "help": ""
  },
  {
    "name": "GetOutput",
    "description": "",
    "args": [
      "id string"
    ],
    "return": [
      "*Port"
    ],
    "help": ""
  },
  {
    "name": "OutputExists",
    "description": "",
    "args": [
      "id string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "SetOutput",
    "description": "",
    "args": [
      "portID string",
      "value any"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "PublishValue",
    "description": "PublishValue eventbus",
    "args": [
      "portID string"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "PublishCommand",
    "description": "",
    "args": [
      "command *ExtendedCommand"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "Subscribe",
    "description": "",
    "args": [
      "topic string",
      "handlerID string",
      "callBack "
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SubscribePayload",
    "description": "",
    "args": [
      "topic string",
      "handlerID string",
      "opts *EventbusOpts",
      "callBack "
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetRootObject",
    "description": "GetRootObject Obj tree",
    "args": [
      "uuid string"
    ],
    "return": [
      "Object",
      "error"
    ],
    "help": ""
  },
  {
    "name": "PrintObjectTree",
    "description": "",
    "args": [
      "objects "
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetCompleteChain",
    "description": "",
    "args": [
      "objects ",
      "uuid string"
    ],
    "return": [
      "Chain"
    ],
    "help": ""
  },
  {
    "name": "RunValidation",
    "description": "ValidationBuilder validation for example, you want to add a new network so lets run some checks eg; is network interface available",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "AddValidation",
    "description": "",
    "args": [
      "key string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "DeleteValidation",
    "description": "",
    "args": [
      "key string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetValidations",
    "description": "",
    "args": null,
    "return": [
      ""
    ],
    "help": ""
  },
  {
    "name": "GetValidation",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "*ErrorsAndValidation",
      "bool"
    ],
    "help": ""
  },
  {
    "name": "SetError",
    "description": "",
    "args": [
      "key string",
      "err error"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SetValidationError",
    "description": "",
    "args": [
      "key string",
      "m *ValidationMessage"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SetHalt",
    "description": "",
    "args": [
      "key string",
      "m *ValidationMessage"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SetValidation",
    "description": "",
    "args": [
      "key string",
      "m *ValidationMessage"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SetStatus",
    "description": "SetStatus -------------------STATS INFO------------------\nthis is for the node status",
    "args": [
      "status ObjectStatus"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "SetLoopCount",
    "description": "",
    "args": [
      "count uint"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetLoopCount",
    "description": "",
    "args": null,
    "return": [
      "uint"
    ],
    "help": ""
  },
  {
    "name": "IncrementLoopCount",
    "description": "",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "ResetLoopCount",
    "description": "",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "GetStats",
    "description": "",
    "args": null,
    "return": [
      "*runtime.ObjectStats"
    ],
    "help": ""
  },
  {
    "name": "SetInfo",
    "description": "SetInfo -------------------OBJECT INFO------------------",
    "args": [
      "info *runtime.Info"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetInfo",
    "description": "",
    "args": null,
    "return": [
      "*runtime.Info"
    ],
    "help": ""
  },
  {
    "name": "GetID",
    "description": "id",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetObjectType",
    "description": "GetObjectType type is for example a driver, service, logic",
    "args": null,
    "return": [
      "ObjectType"
    ],
    "help": ""
  },
  {
    "name": "GetUUID",
    "description": "uuid, set from Meta",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetName",
    "description": "name, set from Meta",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "SetName",
    "description": "",
    "args": [
      "v string"
    ],
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetCategory",
    "description": "category",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetWorkingGroup",
    "description": "working group; a group of objects that work together like a network driver",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetWorkingGroupParent",
    "description": "",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetWorkingGroupLeader",
    "description": "",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetWorkingGroupLeaderObjectUUID",
    "description": "",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetPluginName",
    "description": "plugin",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetMustLiveInObjectType",
    "description": "GetMustLiveInObjectType these are needed to know where a know will site in the sidebar in the UI",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "GetMustLiveParent",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "GetRequiresLogger",
    "description": "",
    "args": null,
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "AddLogger",
    "description": "",
    "args": [
      "trace *Logger"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "Logger",
    "description": "",
    "args": null,
    "return": [
      "*Logger",
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetLoggerInfo",
    "description": "",
    "args": null,
    "return": [
      "[]string",
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetSchema",
    "description": "scheam",
    "args": null,
    "return": [
      "*schema.Generated"
    ],
    "help": ""
  },
  {
    "name": "GetSettings",
    "description": "settings",
    "args": null,
    "return": [
      "*runtime.Settings"
    ],
    "help": ""
  },
  {
    "name": "SetSettings",
    "description": "",
    "args": [
      "settings *runtime.Settings"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetMeta",
    "description": "GetMeta  meta will also set the Obj-name at parentUUID",
    "args": null,
    "return": [
      "*runtime.Meta"
    ],
    "help": ""
  },
  {
    "name": "SetMeta",
    "description": "",
    "args": [
      "meta *runtime.Meta"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetPermissions",
    "description": "permissions",
    "args": null,
    "return": [
      "*runtime.Permissions"
    ],
    "help": ""
  },
  {
    "name": "GetRequirements",
    "description": "requirements",
    "args": null,
    "return": [
      "*runtime.Requirements"
    ],
    "help": ""
  },
  {
    "name": "AddTag",
    "description": "tags",
    "args": [
      "tag string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "AddTags",
    "description": "",
    "args": [
      "tags ...string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetTag",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "HasTag",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "GetTags",
    "description": "",
    "args": null,
    "return": [
      "[]string"
    ],
    "help": ""
  },
  {
    "name": "AddMetaTags",
    "description": "",
    "args": [
      "key string",
      "value string"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "GetMetaTag",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetMetaTags",
    "description": "",
    "args": null,
    "return": [
      ""
    ],
    "help": ""
  },
  {
    "name": "HasMetaTag",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "HasMetaTagValue",
    "description": "",
    "args": [
      "key string",
      "value string"
    ],
    "return": [
      "bool"
    ],
    "help": ""
  },
  {
    "name": "SetCache",
    "description": "",
    "args": [
      "key string",
      "data any",
      "expiration time.Duration",
      "overwriteExisting bool"
    ],
    "return": [
      "error"
    ],
    "help": ""
  },
  {
    "name": "GetCache",
    "description": "",
    "args": [
      "key string"
    ],
    "return": [
      "any",
      "bool"
    ],
    "help": ""
  },
  {
    "name": "CacheAll",
    "description": "",
    "args": null,
    "return": [
      ""
    ],
    "help": ""
  }
]

`
