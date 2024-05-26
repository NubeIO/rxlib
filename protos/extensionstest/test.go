package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/NubeIO/rxlib/ginlib"
	"github.com/NubeIO/rxlib/protos/extensionstest/news"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

var ctx context.Context = context.Background()

type extension struct {
	name       string
	server     *ginlib.Server
	grpcClient runtime.RuntimeServiceClient
	bootTime   string
	pallet     []reactive.Object
	runtime    []reactive.Object
}

type extensionInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	BootTime string `json:"bootTime"`
}

// registerExtension register the extension with the server
func (inst *extension) registerExtension() error {
	s := inst.grpcClient
	regCtx, regCancel := context.WithTimeout(context.Background(), time.Second)
	defer regCancel()
	info := &runtime.Extension{Name: "ExampleExtension", Uuid: inst.name, Pallet: reactive.ConvertObjects(inst.pallet)}
	_, err := s.RegisterExtension(regCtx, info)
	if err != nil {
		return fmt.Errorf("could not register extension: %v", err)
	}
	fmt.Println("Registered extension")
	return nil
}

func (inst *extension) bootServer(opts *ginlib.Opts) {
	inst.server = ginlib.NewServer(opts)
}

func (inst *extension) infoRoute() {
	inst.server.AddGetRoute("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, extensionInfo{
			Name:     inst.name,
			Version:  "1.0.0",
			BootTime: inst.bootTime,
		})
	})
}

func (inst *extension) pingRoute() {
	inst.server.AddGetRoute("/api/ping", func(c *gin.Context) {
		r, err := inst.grpcClient.Ping(ctx, &runtime.PingRequest{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to ping server"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("message: %s port: %d", r.Message, inst.server.GetPort())})
	})
}

func (inst *extension) palletRoute() {
	inst.server.AddGetRoute("/api/pallet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pallet": reactive.ConvertObjects(inst.pallet)})
	})
}

func (inst *extension) bootGRPC() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	inst.grpcClient = runtime.NewRuntimeServiceClient(conn)
}

func (inst *extension) buildPallet() {
	baseObj := reactive.New("news", nil)
	instance := news.New(nil)
	obj := instance.New(baseObj)
	obj.GetInfo().PluginName = inst.name
	inst.pallet = append(inst.pallet, obj)
}

func main() {
	c := &extension{
		name:     "test",
		bootTime: time.Now().String(),
	}
	c.buildPallet()
	port := flag.String("port", "4000", "Port number for the server")
	flag.Parse()
	c.bootGRPC() // Ensure gRPC is connected before setting up the server
	c.bootServer(&ginlib.Opts{
		Port: *port,
	})
	c.infoRoute()
	c.pingRoute()
	c.palletRoute()
	c.registerExtension()
	c.server.Run()
}
