package ndocs

var RuntimesString = `
[
  {
    "name": "Get",
    "description": "Get get all objects []Object",
    "args": null,
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "AddObjects",
    "description": "AddObjects add object to runtime",
    "args": null,
    "return": null,
    "help": ""
  },
  {
    "name": "ToObjectsConfig",
    "description": "ToObjectsConfig convert to ObjectConfig, used when needed as JSON",
    "args": [
      "objects []Object"
    ],
    "return": [
      "[]*runtime.ObjectConfig"
    ],
    "help": ""
  },
  {
    "name": "GetObjectsUUIDs",
    "description": "",
    "args": [
      "objects []Object"
    ],
    "return": [
      "[]string"
    ],
    "help": ""
  },
  {
    "name": "GetObjectsConfig",
    "description": "",
    "args": null,
    "return": [
      "[]*runtime.ObjectConfig"
    ],
    "help": ""
  },
  {
    "name": "GetObjectConfig",
    "description": "",
    "args": [
      "uuid string"
    ],
    "return": [
      "*runtime.ObjectConfig"
    ],
    "help": ""
  },
  {
    "name": "GetObjectConfigByID",
    "description": "",
    "args": [
      "objectID string"
    ],
    "return": [
      "*runtime.ObjectConfig"
    ],
    "help": ""
  },
  {
    "name": "GetObjectValues",
    "description": "",
    "args": [
      "objectUUID string"
    ],
    "return": [
      "[]*runtime.PortValue"
    ],
    "help": ""
  },
  {
    "name": "GetObjectsValues",
    "description": "",
    "args": null,
    "return": [
      ""
    ],
    "help": ""
  },
  {
    "name": "GetObjectsValuesPaginate",
    "description": "",
    "args": [
      "parentUUID string",
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectValuesPagination"
    ],
    "help": ""
  },
  {
    "name": "ObjectsPagination",
    "description": "",
    "args": [
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectPagination"
    ],
    "help": ""
  },
  {
    "name": "PaginateGetAllByID",
    "description": "",
    "args": [
      "objectID string",
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectPagination"
    ],
    "help": ""
  },
  {
    "name": "PaginateGetChildObjects",
    "description": "",
    "args": [
      "parentUUID string",
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectPagination"
    ],
    "help": ""
  },
  {
    "name": "PaginateGetAllByName",
    "description": "",
    "args": [
      "name string",
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectPagination"
    ],
    "help": ""
  },
  {
    "name": "PaginateGetChildObjectsByWorkingGroup",
    "description": "",
    "args": [
      "objectUUID string",
      "workingGroup string",
      "pageNumber int",
      "pageSize int"
    ],
    "return": [
      "*ObjectPagination"
    ],
    "help": ""
  },
  {
    "name": "Delete",
    "description": "",
    "args": null,
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "GetByUUID",
    "description": "",
    "args": [
      "uuid string"
    ],
    "return": [
      "Object"
    ],
    "help": ""
  },
  {
    "name": "GetFirstByID",
    "description": "",
    "args": [
      "objectID string"
    ],
    "return": [
      "Object"
    ],
    "help": ""
  },
  {
    "name": "GetAllByID",
    "description": "",
    "args": [
      "objectID string"
    ],
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "GetFirstByName",
    "description": "",
    "args": [
      "name string"
    ],
    "return": [
      "Object"
    ],
    "help": ""
  },
  {
    "name": "GetAllByName",
    "description": "",
    "args": [
      "name string"
    ],
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "GetChildObjectsByWorkingGroup",
    "description": "",
    "args": [
      "objectUUID string",
      "workingGroup string"
    ],
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "GetChildObjects",
    "description": "",
    "args": [
      "parentUUID string"
    ],
    "return": [
      "[]Object"
    ],
    "help": ""
  },
  {
    "name": "GetAllObjectValues",
    "description": "",
    "args": null,
    "return": [
      "[]*ObjectValue"
    ],
    "help": ""
  },
  {
    "name": "AddObject",
    "description": "",
    "args": [
      "object Object"
    ],
    "return": null,
    "help": ""
  },
  {
    "name": "Command",
    "description": "",
    "args": [
      "cmd *ExtendedCommand"
    ],
    "return": [
      "*runtime.CommandResponse"
    ],
    "help": ""
  },
  {
    "name": "CommandObject",
    "description": "",
    "args": [
      "cmd *ExtendedCommand"
    ],
    "return": [
      "*CommandResponse"
    ],
    "help": ""
  },
  {
    "name": "GetTreeMapRoot",
    "description": "",
    "args": null,
    "return": [
      "*runtime.ObjectsRootMap"
    ],
    "help": ""
  },
  {
    "name": "GetAncestorTreeByUUID",
    "description": "",
    "args": [
      "objectUUID string"
    ],
    "return": [
      "*AncestorTreeNode"
    ],
    "help": ""
  },
  {
    "name": "GetTreeChilds",
    "description": "",
    "args": [
      "objectUUID string"
    ],
    "return": [
      "*AncestorTreeNode"
    ],
    "help": ""
  },
  {
    "name": "AllPlugins",
    "description": "",
    "args": null,
    "return": [
      "[]*plugins.Export"
    ],
    "help": ""
  },
  {
    "name": "GetObjectsPallet",
    "description": "",
    "args": null,
    "return": [
      "*PalletTree"
    ],
    "help": ""
  },
  {
    "name": "Cron",
    "description": "",
    "args": null,
    "return": [
      "scheduler.Cron"
    ],
    "help": ""
  },
  {
    "name": "ExprWithError",
    "description": "",
    "args": [
      "query string"
    ],
    "return": [
      "any",
      "error"
    ],
    "help": ""
  },
  {
    "name": "Expr",
    "description": "",
    "args": [
      "query string"
    ],
    "return": [
      "any"
    ],
    "help": ""
  },
  {
    "name": "System",
    "description": "System os/host system info",
    "args": null,
    "return": [
      "systeminfo.System"
    ],
    "help": ""
  },
  {
    "name": "SystemInfo",
    "description": "",
    "args": null,
    "return": [
      "*systeminfo.Info"
    ],
    "help": ""
  },
  {
    "name": "HistoryManager",
    "description": "",
    "args": null,
    "return": [
      "history.Manager"
    ],
    "help": ""
  },
  {
    "name": "ToStringArray",
    "description": "ToStringArray Conversions",
    "args": [
      "interfaces "
    ],
    "return": [
      "[]string"
    ],
    "help": ""
  },
  {
    "name": "DB",
    "description": "",
    "args": null,
    "return": [
      "pglib.PG"
    ],
    "help": ""
  },
  {
    "name": "Rest",
    "description": "",
    "args": null,
    "return": [
      "restc.Rest"
    ],
    "help": ""
  },
  {
    "name": "Publish",
    "description": "",
    "args": [
      "topic string",
      "body "
    ],
    "return": [
      "string"
    ],
    "help": ""
  },
  {
    "name": "Iam",
    "description": "",
    "args": [
      "rangeStart int",
      "finish int"
    ],
    "return": [
      "Object"
    ],
    "help": ""
  }
]

`
