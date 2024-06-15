package watcher

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/protos/extensionlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"sync"
)

type Instance struct {
	reactive.Object
	stop          chan struct{}
	once          sync.Once
	outputUpdated func(message *runtime.Command)
	filepath      string
	watcher       *fsnotify.Watcher
}

var infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func New(outputUpdated func(message *runtime.Command)) extensionlib.PluginObject {
	obj := new(Instance)
	obj.outputUpdated = outputUpdated
	return obj
}

func (inst *Instance) New(object reactive.Object, opts ...any) reactive.Object {
	info := rxlib.NewObjectInfo().
		SetID("watcher").
		SetPluginName("ext-math").
		SetCategory("util").
		SetCallResetOnDeploy().
		SetObjectType(rxlib.Service).
		SetAllPermissions().
		Build()

	object.SetInfo(info)
	object.NewOutputPort(&runtime.Port{
		Id:        "output",
		Name:      "output",
		Direction: string(rxlib.Output),
		DataType:  priority.TypeJSON,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "file",
		Name:      "file",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeString,
	})
	inst.Object = object
	inst.stop = make(chan struct{})
	return inst
}

func (inst *Instance) OutputUpdated(message *runtime.Command) {
	inst.outputUpdated(message)
}

func (inst *Instance) Start() error {
	log.Println("started filewatcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return err
	}
	inst.watcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Check if the file was modified
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					jsonData, err := json.Marshal(event)
					if err == nil {
						json := string(jsonData)
						inst.OutputUpdated(&runtime.Command{
							Key:              "update-outputs",
							TargetObjectUUID: inst.GetMeta().GetObjectUUID(),
							PortValues: []*runtime.PortValue{&runtime.PortValue{
								PortID:    "output",
								JsonValue: json,
								DataType:  priority.TypeJSON,
							}},
						})
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	inst.filepath = "/home/manny/nebu/rxlib/protos/filewatcher/index.txt"
	err = watcher.Add(inst.filepath)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (inst *Instance) Reset() error {
	inst.Delete()
	inst.Start()
	return nil
}

func (inst *Instance) Delete() error {
	inst.once.Do(func() {
		close(inst.stop)
	})
	inst.watcher.Close()
	return nil
}

func (inst *Instance) Handler(p *runtime.MessageRequest) {
	infoLog.Println("filewatcher Handler")
	if p == nil {
		return
	}
	cmd := p.GetCommand()
	if cmd == nil {
		return
	}

	for _, value := range cmd.GetPortValues() {
		for _, d := range value.PortIDs {
			if d == "file" {
				file := value.StringValue
				if inst.filepath == file {
					return
				}
				inst.Delete()
				inst.filepath = file
			}
		}
	}

	log.Println("started filewatcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return
	}
	inst.watcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Check if the file was modified
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					jsonData, err := json.Marshal(event)
					if err == nil {
						json := string(jsonData)
						inst.OutputUpdated(&runtime.Command{
							Key:              "update-outputs",
							TargetObjectUUID: inst.GetMeta().GetObjectUUID(),
							PortValues: []*runtime.PortValue{&runtime.PortValue{
								PortID:    "output",
								JsonValue: json,
								DataType:  priority.TypeJSON,
							}},
						})
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(inst.filepath)
	if err != nil {
		log.Println(err)
	}
}
