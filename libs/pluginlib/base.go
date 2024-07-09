package pluginlib

import (
	"flag"
	"fmt"
	"github.com/NubeIO/rxlib/ginlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"google.golang.org/grpc"
	"log"
	"time"
)

var messages map[string]interface{}

type GeneratePlugin func(outputUpdated func(message *runtime.Command)) PluginObject

type PluginObject interface {
	New(object reactive.Object, opts ...any) reactive.Object
	OutputUpdated(message *runtime.Command)
	Start() error
}

type Plugins struct {
	name       string
	server     *ginlib.Server
	grpcClient runtime.RuntimeServiceClient
	bootTime   string
	pallet     map[string]reactive.Object
	registry   map[string]GeneratePlugin
	runtime    []reactive.Object
	callbacks  map[string]func(message *runtime.MessageRequest)
	stream     runtime.RuntimeService_PluginStreamClient
	port       string
	grpcPort   string
}

type extensionInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	BootTime string `json:"bootTime"`
}

func New(name string) *Plugins {
	cli := &Plugins{
		name:     name,
		bootTime: time.Now().String(),
	}
	messages = make(map[string]interface{})
	cli.pallet = make(map[string]reactive.Object)
	cli.registry = make(map[string]GeneratePlugin)
	cli.runtime = []reactive.Object{}
	cli.callbacks = map[string]func(message *runtime.MessageRequest){}

	port := flag.String("port", "4000", "Port number for the server")
	grpcPort := flag.String("grpc", "9090", "Port number for grpc server")
	flag.Parse()
	cli.port = *port
	cli.grpcPort = *grpcPort

	cli.bootGRPC() // Ensure gRPC is connected before setting up the server
	cli.bootServer(&ginlib.Opts{
		Port: cli.port,
	})
	cli.infoRoute()
	cli.pingRoute()
	cli.runtimeRoute()
	cli.palletRoute()
	cli.messagesRoute()
	cli.streamRoute()

	return cli
}

func (inst *Plugins) bootGRPC() {
	target := fmt.Sprintf("localhost:%s", inst.grpcPort)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	inst.grpcClient = runtime.NewRuntimeServiceClient(conn)
}
