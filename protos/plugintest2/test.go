package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := runtime.NewRuntimeServiceClient(conn)

	// Register the plugin with a separate context for registration
	regCtx, regCancel := context.WithTimeout(context.Background(), time.Second)
	defer regCancel()
	info := &runtime.PluginInfo{Name: "ExamplePlugin", Uuid: "12345", Address: "localhost:9091"}
	_, err = c.RegisterPlugin(regCtx, info)
	if err != nil {
		log.Fatalf("could not register plugin: %v", err)
	}
	fmt.Println("Registered plugin")

	// Start bidirectional streaming with a different context
	streamCtx, streamCancel := context.WithCancel(context.Background())
	defer streamCancel()
	stream, err := c.PluginStreamMessages(streamCtx)
	if err != nil {
		log.Fatalf("failed to start streaming: %v", err)
	}

	// Send a message to the server
	if err := stream.Send(&runtime.Message{Uuid: info.Uuid, Content: "Hello from plugin"}); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	// Receive messages from the server
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}
		fmt.Printf("Received message from server: %s\n", in.Content)
	}
}
