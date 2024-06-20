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
	"os/user"
	"sync"
	"time"
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
		SetPluginName("test").
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
		DataType:  priority.TypeString,
	})
	object.NewInputPort(&runtime.Port{
		Id:        "file",
		Name:      "file",
		Direction: string(rxlib.Input),
		DataType:  priority.TypeString,
	})
	inst.Object = object
	inst.stop = make(chan struct{})
	fmt.Printf("init new file watch: %s \n", inst.GetMeta().GetObjectUUID())
	return inst
}

func (inst *Instance) OutputUpdated(message *runtime.Command) {
	inst.outputUpdated(message)
}

func (inst *Instance) startWatcher(filePath string) {
	log.Println("started filewatcher")
	interval := time.Second // Polling interval

	// Get initial file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	lastModTime := fileInfo.ModTime()

	for {
		time.Sleep(interval)

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatalf("Failed to get file info: %v", err)
		}

		currentModTime := fileInfo.ModTime()
		if currentModTime != lastModTime {
			change := fmt.Sprintf("File %s changed at %v", filePath, currentModTime)
			lastModTime = currentModTime
			if err == nil {
				inst.OutputUpdated(&runtime.Command{
					Key:              "update-outputs",
					TargetObjectUUID: inst.GetMeta().GetObjectUUID(),
					PortValues: []*runtime.PortValue{&runtime.PortValue{
						PortID:      "output",
						StringValue: change,
						DataType:    priority.TypeString,
					}},
				})
			}
		}
	}
}

func (inst *Instance) Start() error {
	log.Println("started filewatcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return err
	}
	inst.watcher = watcher
	inst.filepath = fmt.Sprintf("%s/test.txt", getHome())
	log.Printf("file to watch: %s\n", inst.filepath)
	err = watcher.Add(inst.filepath)
	if err != nil {
		log.Println(err)
	}
	inst.startWatcher(inst.filepath)
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
	fmt.Println("filewatcher Handler !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1")
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
								DataType:  priority.TypeString,
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

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir
}
