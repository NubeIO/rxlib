package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/NubeIO/rxlib/ginlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

var ctx context.Context = context.Background()

type client struct {
	server     *ginlib.Server
	grpcClient runtime.RuntimeServiceClient
}

func (inst *client) bootServer(opts *ginlib.Opts) {
	inst.server = ginlib.NewServer(opts)
}

func (inst *client) ping() {
	inst.server.AddGetRoute("/api/ping", func(c *gin.Context) {
		r, err := inst.grpcClient.Ping(ctx, &runtime.PingRequest{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to ping server"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("message: %s port: %d", r.Message, inst.server.GetPort())})
	})
}

func (inst *client) bootGRPC() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	inst.grpcClient = runtime.NewRuntimeServiceClient(conn)
}

func main() {
	c := &client{}
	port := flag.String("port", "4000", "Port number for the server")
	flag.Parse()
	c.bootGRPC() // Ensure gRPC is connected before setting up the server
	c.bootServer(&ginlib.Opts{
		Port: *port,
	})
	c.ping()
	c.server.Run()
}
