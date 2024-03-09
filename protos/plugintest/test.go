package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/rxlib/protos/plugintest/counter"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type client struct {
	callbacks map[string]func(message *runtime.MessageRequest)
	stream    runtime.RuntimeService_PluginStreamMessagesClient
	pallet    []reactive.Object
	runtime   []reactive.Object
}

func (cli *client) sendMessage(content string) error {
	if cli.stream == nil {
		return fmt.Errorf("stream is not initialized")
	}

	// Send a message to the server
	if err := cli.stream.Send(&runtime.MessageRequest{Uuid: "1234"}); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (cli *client) connectWithRetry() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for {
		conn, err = grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("could not connect: %v", err)
			time.Sleep(30 * time.Second) // Retry after 30 seconds
			continue
		}
		break
	}

	return conn, nil
}

func (cli *client) registerPlugin(conn *grpc.ClientConn) error {
	c := runtime.NewRuntimeServiceClient(conn)

	// Register the plugin with a separate context for registration
	regCtx, regCancel := context.WithTimeout(context.Background(), time.Second)
	defer regCancel()
	info := &runtime.PluginInfo{Name: "ExamplePlugin", Uuid: "1234", Pallet: reactive.ConvertObjects(cli.pallet)}
	_, err := c.RegisterPlugin(regCtx, info)
	if err != nil {
		return fmt.Errorf("could not register plugin: %v", err)
	}
	fmt.Println("Registered plugin")
	return nil
}

func (cli *client) startStreaming(ctx context.Context, conn *grpc.ClientConn) error {
	c := runtime.NewRuntimeServiceClient(conn)

	// Start bidirectional streaming with the given context
	stream, err := c.PluginStreamMessages(ctx)
	if err != nil {
		return fmt.Errorf("failed to start streaming: %v", err)
	}
	cli.stream = stream

	// Send a message to the server
	if err := stream.Send(&runtime.MessageRequest{Uuid: "1234"}); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	// Receive messages from the server
	for {
		select {
		case <-ctx.Done():
			return nil // Context canceled, exit the loop
		default:
			in, err := stream.Recv()
			if err == io.EOF {
				return nil // Server closed the stream
			}
			if err != nil {
				return fmt.Errorf("failed to receive message: %v", err)
			}
			if callback, ok := cli.callbacks[in.Topic]; ok {
				callback(in)
			} else {
				fmt.Printf("Received message from server unknown: %s\n", in.Topic)
			}
		}
	}
}

func (cli *client) callbackOne(message *runtime.MessageRequest) {
	// Send a message to the server
	if err := cli.sendMessage("response for 1"); err != nil {
		fmt.Printf("failed to send message: %v\n", err)
	}
}

func (cli *client) callbackTwo(message *runtime.MessageRequest) {
	// Send a message to the server
	if err := cli.sendMessage("response for 2"); err != nil {
		fmt.Printf("failed to send message: %v\n", err)
	}
}

func (cli *client) getPallet() {
	baseObj := reactive.New("counter-2")
	instance := counter.New()
	obj := instance.New(baseObj)
	cli.pallet = append(cli.pallet, obj)
}

func (cli *client) counter() reactive.Object {
	baseObj := reactive.New("counter-2")
	instance := counter.New()
	base := instance.New(baseObj)
	return base
}

func main() {
	cli := &client{
		callbacks: make(map[string]func(message *runtime.MessageRequest)),
		pallet:    []reactive.Object{},
		runtime:   []reactive.Object{},
	}
	cli.getPallet()

	var err error

	instance := cli.counter()

	cli.runtime = append(cli.runtime, instance)

	cli.callbacks["one"] = instance.Handler

	var conn *grpc.ClientConn

	for {
		conn, err = cli.connectWithRetry()
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()

		if err := cli.registerPlugin(conn); err != nil {
			log.Fatalf("could not register plugin: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := cli.startStreaming(ctx, conn); err != nil {
			log.Printf("streaming error: %v, reconnecting...", err)
			continue
		}
	}
}
