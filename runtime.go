package rxlib

import (
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib/config"
	"github.com/NubeIO/rxlib/libs/alarm"
	"github.com/NubeIO/rxlib/libs/chat"
	"github.com/NubeIO/rxlib/libs/history"
	"github.com/NubeIO/rxlib/libs/jsonutils"
	"github.com/NubeIO/rxlib/libs/pglib"
	"github.com/NubeIO/rxlib/libs/restc"
	"github.com/NubeIO/rxlib/libs/schedules"
	systeminfo "github.com/NubeIO/rxlib/libs/system"
	"github.com/NubeIO/rxlib/plugins"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/NubeIO/scheduler"
	"log"
	"sync"
)

type Runtime interface {
	// Get all objects []Object
	Get() []Object
	// AddObjects adds objects to runtime
	AddObjects([]Object)
	// AddObject adds an object
	AddObject(object Object)
	// Deploy deploy a flow
	Deploy(body *Deploy) *DeployResponse
	// ToObjectConfig converts to ObjectConfig, used when needed as JSON
	ToObjectConfig(objects Object) *runtime.ObjectConfig
	// ToObjectsConfig converts to ObjectConfig, used when needed as JSON
	ToObjectsConfig(objects []Object) []*runtime.ObjectConfig
	// GetObjectsUUIDs returns UUIDs of objects
	GetObjectsUUIDs(objects []Object) []string
	// GetObjectsRootConfig returns config of all objects for the root
	GetObjectsRootConfig() []*runtime.ObjectConfig
	// GetObjectsConfig returns config of all objects
	GetObjectsConfig() []*runtime.ObjectConfig
	// GetObjectConfig returns config of an object by UUID
	GetObjectConfig(uuid string) *runtime.ObjectConfig
	// GetObjectConfigByID returns config of an object by ID
	GetObjectConfigByID(objectID string) *runtime.ObjectConfig

	// GetObjectValues returns values of an object by UUID
	GetObjectValues(objectUUID string) []*runtime.PortValue
	// GetObjectsValues returns values of all objects, if the parentUUID is passed in it will return all the child object port values. FYI will not return the parent value
	GetObjectsValues(parentUUID string) []*runtime.PortValue
	// GetObjectsValuesPaginate returns paginated values of objects
	GetObjectsValuesPaginate(parentUUID string, pageNumber, pageSize int) *ObjectValuesPagination

	// ObjectsPagination paginates objects
	ObjectsPagination(pageNumber, pageSize int) *ObjectPagination
	// PaginateGetAllByID paginates objects by ID
	PaginateGetAllByID(objectID string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetChildObjects paginates child objects
	PaginateGetChildObjects(parentUUID string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetAllByName paginates objects by name
	PaginateGetAllByName(name string, pageNumber, pageSize int) *ObjectPagination
	// PaginateGetChildObjectsByWorkingGroup paginates child objects by working group
	PaginateGetChildObjectsByWorkingGroup(objectUUID, workingGroup string, pageNumber, pageSize int) *ObjectPagination

	// Delete deletes runtime
	Delete() string
	// DeleteByUUID deletes runtime by UUID
	DeleteByUUID(uuid string) error
	// GetByUUID gets object by UUID; eg GetByUUID("abc").GetName(), GetByUUID("abc").GetInputs()
	GetByUUID(uuid string) Object
	// GetFirstByID gets first object by ID eg GetFirstByID("abc").GetID(), GetFirstByID("abc").GetTags()
	GetFirstByID(objectID string) Object
	// GetAllByID gets all objects by ID
	GetAllByID(objectID string) []Object
	// GetFirstByName gets first object by name
	GetFirstByName(name string) Object
	// GetAllByName gets all objects by name
	GetAllByName(name string) []Object

	// GetChildObjectsByWorkingGroup gets child objects by working group
	GetChildObjectsByWorkingGroup(objectUUID, workingGroup string) []Object
	// GetChildObjects gets child objects
	GetChildObjects(parentUUID string) []Object
	// GetAllObjectValues gets all object values
	GetAllObjectValues() []*ObjectValue
	// Command executes a command
	Command(cmd *ExtendedCommand) *runtime.CommandResponse
	// CommandObject executes a command for an object
	CommandObject(cmd *ExtendedCommand) *CommandResponse

	// GetTreeMapRoot gets the root of the object tree map
	GetTreeMapRoot() *runtime.ObjectsRootMap
	// GetAncestorTreeByUUID gets ancestor tree by UUID
	GetAncestorTreeByUUID(objectUUID string) *runtime.AncestorObjectTree
	// GetTreeChilds gets child nodes of an object
	GetTreeChilds(objectUUID string) *runtime.AncestorObjectTree

	// AllPlugins returns all plugins
	AllPlugins() []*plugins.Export

	// GetObjectsPallet returns objects pallet tree
	GetObjectsPallet() *PalletTree

	// ExprWithError run a system query and returns an error. eg; Expr("filter(objects, .GetID() == "rubix-manager""))  see docs https://github.com/expr-lang/expr
	ExprWithError(query string) (any, error)

	// Expr run a system query. eg; Expr("filter(objects, .GetID() == "rubix-manager""))  see docs https://github.com/expr-lang/expr
	Expr(query string) any

	// System get host info, networking, memory and stats. eg; System().GetIP()
	System() systeminfo.System

	// SystemInfo return the host machine info eg; update and timezone
	SystemInfo() *systeminfo.Info

	// HistoryManager get ros history manager. eg; HistoryManager().AllHistories()
	HistoryManager() history.Manager

	// Cron gets cron/job scheduler  Cron().All()
	Cron() scheduler.Scheduler

	// ScheduleManager weekly and exception time scheduler eg; ScheduleManager().All()
	ScheduleManager() schedules.Manager

	// Rest is a rest client for making HTTP requests, Example Rest().Execute("GET", "http://localhost:8080/api")
	Rest() restc.Rest

	// DB access the postgres database; eg DG().Select("select *")
	DB() pglib.PG

	// AlarmManager usage is through the alarm-manager object
	AlarmManager() alarm.Manager

	// Publish a mqtt message
	Publish(topic string, body interface{}) (err string)
	// RequestResponse a mqtt message want wait for a repose
	RequestResponse(timeoutSeconds int, publishTopic, responseTopic, requestUUID string, body interface{}) *mqttwrapper.Response
	// Client rubix-os mqtt client; eg Client.RQL()
	Client() ROSClient

	Iam(rangeStart, finish int) Object
	IamConfig(rangeStart, rangeEnd int) *runtime.ObjectConfig

	// ObjectSync sync all the objects to the postgres db;
	ObjectSync(forceSync bool, opts *SyncOptions) error
	// HistorySync sync all the histories to the postgres db;
	HistorySync(forceSync bool, opts *SyncOptions) error

	// UUID generates a UUID
	UUID() string

	// ObjectBuilder is used for build an object. Use case if for building an object using RQL and then deploying. eg; ObjectBuilder({"objectID: trigger"}).ToObject()
	ObjectBuilder(body *Builder) *Builder

	// ToStringArray Conversions
	ToStringArray(interfaces interface{}) []string

	// JSON helper functions to work with JSON, you can also use gjson see docs; https://github.com/tidwall/gjson
	JSON() jsonutils.JSON

	// ChatGPT Send a message to chatGPT, if the model is empty it will use 3.5. The preloaded data help if you want to explain to chatGPT some extra info to help with the users query
	ChatGPT(token, body, preloaded string, model ...string) *chat.Response
	ChatBot(token, body string, model ...string) *chat.Response
	Config() *config.Configuration

	// Manager get the rubix manager
	Manager() Object
}

type RuntimeSettings struct {
	GlobalUUID string
	GlobalID   string
	RootDir    string
	Version    string
}

type RuntimeOpts struct {
	Scheduler       scheduler.Scheduler
	MQTTClient      mqttwrapper.MQTT
	RuntimeSettings *RuntimeSettings
}

func NewRuntime(objs []Object, opts *RuntimeOpts) Runtime {
	r := &RuntimeImpl{
		tree:       &tree{},
		mqttClient: opts.MQTTClient,
	}
	r.objects = objs
	r.scheduler = opts.Scheduler
	r.hist = history.NewHistoryManager("ros")
	connString := "postgresql://postgres:postgres@localhost/postgres"
	db, err := pglib.New(connString)
	if err != nil {
		fmt.Println("runtime init DB: ", err)
	}
	r.db = db
	r.rest = restc.New()
	r.alarmManager = alarm.NewAlarmManager("runtime")
	r.scheduleManager = schedules.New()
	r.jsonUtils = jsonutils.New()
	r.runtimeSettings = opts.RuntimeSettings
	if r.mqttClient == nil {
		log.Fatal("Runtime() mqtt client can not be empty")
	}
	if r.runtimeSettings == nil {
		log.Fatal("Runtime() runtimeSettings can not be empty")
	}
	r.config = config.Get()
	r.client = NewRosClient(opts.MQTTClient, r.runtimeSettings)
	return r
}

func (inst *RuntimeImpl) Get() []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	return inst.objects
}

type RuntimeImpl struct {
	objects         []Object
	objectsFiltered []Object
	PluginsExport   []*plugins.Export
	err             error // To handle errors in query chain
	where           string
	field           string
	mutex           sync.RWMutex
	response        *CommandResponse
	parsedCommand   *ParsedCommand
	command         *ExtendedCommand
	tree            *tree
	addedObject     bool
	scheduler       scheduler.Scheduler
	hist            history.Manager
	db              pglib.PG
	rest            restc.Rest
	mqttClient      mqttwrapper.MQTT
	alarmManager    alarm.Manager
	scheduleManager schedules.Manager
	jsonUtils       jsonutils.JSON
	runtimeSettings *RuntimeSettings
	client          ROSClient
	config          *config.Configuration
}

func (inst *RuntimeImpl) JSON() jsonutils.JSON {
	return inst.jsonUtils
}

func (inst *RuntimeImpl) Config() *config.Configuration {
	return inst.config
}

func (inst *RuntimeImpl) AlarmManager() alarm.Manager {
	return inst.alarmManager
}

func (inst *RuntimeImpl) Manager() Object {
	return inst.GetFirstByID("rubix-manager")
}

func (inst *RuntimeImpl) ScheduleManager() schedules.Manager {
	return inst.scheduleManager
}

func (inst *RuntimeImpl) Client() ROSClient {
	return inst.client
}

func (inst *RuntimeImpl) Publish(topic string, body interface{}) string {
	if inst.mqttClient == nil {
		return "client is empty"
	}
	err := inst.mqttClient.Publish(topic, body)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (inst *RuntimeImpl) RequestResponse(timeoutSeconds int, publishTopic, responseTopic, requestUUID string, body interface{}) *mqttwrapper.Response {
	if inst.mqttClient == nil {
		return &mqttwrapper.Response{
			Error: "client is empty",
		}
	}
	return inst.mqttClient.RequestResponse(timeoutSeconds, publishTopic, responseTopic, requestUUID, body)
}

func (inst *RuntimeImpl) AddObjects(objects []Object) {
	inst.objects = objects
}

func (inst *RuntimeImpl) HistoryManager() history.Manager {
	return inst.hist
}

func (inst *RuntimeImpl) DB() pglib.PG {
	return inst.db
}

func (inst *RuntimeImpl) Rest() restc.Rest {
	return inst.rest
}

func (inst *RuntimeImpl) System() systeminfo.System {
	return systeminfo.NewSystem()
}

func (inst *RuntimeImpl) SystemInfo() *systeminfo.Info {
	s := inst.System().Info()
	if s == nil {
		return nil
	}
	s.GlobalID = inst.runtimeSettings.GlobalID
	s.Version = inst.runtimeSettings.Version
	return s
}

func (inst *RuntimeImpl) Cron() scheduler.Scheduler {
	return inst.scheduler
}

func (inst *RuntimeImpl) Delete() string {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	c := len(inst.objects)
	inst.objects = nil
	d := len(inst.objects)
	return fmt.Sprintf("count deleted: %d current: %d", c, d)
}

func (inst *RuntimeImpl) DeleteByUUID(uuid string) error {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	deleted := false
	for i, o := range inst.objects {
		if o.GetUUID() == uuid {
			inst.objects = append(inst.objects[:i], inst.objects[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		return fmt.Errorf("not found object with uuid: %s", uuid)
	}
	return nil
}

func (inst *RuntimeImpl) ChatBot(token, body string, model ...string) *chat.Response {
	var m string
	if len(model) > 0 {
		m = model[0]
	}
	return chat.NewMessage(&chat.Chat{
		Token:   token,
		Content: body,
		PreLoad: pglib.ChatGPTInfo(),
		Model:   m,
	})
}

func (inst *RuntimeImpl) ChatGPT(token, body, preloaded string, model ...string) *chat.Response {
	var m string
	if len(model) > 0 {
		m = model[0]
	}
	return chat.NewMessage(&chat.Chat{
		Token:   token,
		Content: body,
		PreLoad: preloaded,
		Model:   m,
	})
}
