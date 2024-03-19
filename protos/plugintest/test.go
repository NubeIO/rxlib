package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/rxlib/libs/unixbus"
	"github.com/NubeIO/rxlib/protos/plugintest/add"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type client struct {
	pluginName       string
	callbacks        map[string]func(message *runtime.MessageRequest)
	stream           runtime.RuntimeService_PluginStreamMessagesClient
	pallet           []reactive.Object
	runtime          []reactive.Object
	serverConnection runtime.RuntimeServiceClient
}

func (cli *client) sendMessage(content string) error {
	if cli.stream == nil {
		return fmt.Errorf("stream is not initialized")
	}

	// Send a message to the server
	if err := cli.stream.Send(&runtime.MessageRequest{Uuid: cli.pluginName}); err != nil {
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

func (cli *client) registerPlugin() error {
	c := cli.serverConnection

	// Register the plugin with a separate context for registration
	regCtx, regCancel := context.WithTimeout(context.Background(), time.Second)
	defer regCancel()
	info := &runtime.PluginInfo{Name: "ExamplePlugin", Uuid: cli.pluginName, Pallet: reactive.ConvertObjects(cli.pallet)}
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
	if err := stream.Send(&runtime.MessageRequest{Uuid: cli.pluginName}); err != nil {
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
			fmt.Println(in.Key)
			//pprint.PrintJSON(in)
			//if callback, ok := cli.callbacks[in.Key]; ok {
			//	callback(in)
			//} else {
			//	fmt.Printf("Received message from server unknown: %s\n", in.Key)
			//}
		}
	}
}

func (cli *client) addObject(message *runtime.MessageRequest) {
	object := message.GetObject()
	if object == nil {
		if err := cli.sendMessage("failed to get object"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
		return
	}
	objectID := object.GetInfo().GetObjectID()
	if objectID == "" {
		if err := cli.sendMessage("objectID is empty"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
		return
	}

	instance := cli.counter(reactive.ConvertObjectConfig(object), cli.outputCallback)
	if instance != nil {
		cli.runtime = append(cli.runtime, instance)
		cli.callbacks[instance.GetMeta().GetObjectUUID()] = instance.Handler
		if err := cli.sendMessage("response for 1"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
		fmt.Printf("objects count: %d\n", len(cli.runtime))

		return
	} else {
		if err := cli.sendMessage("failed to find object in pallet"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
	}
}

func (cli *client) outputCallback(cmd *runtime.Command) {
	//cli.serverConnection.ObjectInvoke(context.Background(), cmd)

	if err := cli.stream.Send(&runtime.MessageRequest{
		Key:     "invoke",
		Uuid:    cli.pluginName,
		Command: cmd,
	}); err != nil {

	}
}

func (cli *client) getPallet() {
	baseObj := reactive.New("add", nil)
	instance := add.New(nil)
	obj := instance.New(baseObj)
	obj.GetInfo().PluginName = cli.pluginName
	cli.pallet = append(cli.pallet, obj)
}

func (cli *client) existsInPallet(objectID string) bool {
	for _, object := range cli.pallet {
		if object.GetInfo().GetObjectID() == objectID {
			return true
		}
	}
	return false
}

func (cli *client) counter(obj *reactive.BaseObject, outputUpdated func(message *runtime.Command)) reactive.Object {
	instance := add.New(outputUpdated)
	base := instance.New(obj)
	return base
}

const (
	createObject = "create-object"
)

func (cli *client) bus() {
	eventBus := unixbus.NewUnixEventBus("user.topic.123")
	eventBus.Subscribe(func(data interface{}) {
		//marshal, err := json.Marshal(data)
		//if err != nil {
		//	fmt.Println(111, err)
		//	return
		//}
		//var message *runtime.MessageRequest
		//err = json.Unmarshal(marshal, &message)
		//if err != nil {
		//	fmt.Println(222, err)
		//	return
		//}
		//fmt.Println(string(marshal))
	})

}

func main() {
	cli := &client{
		pluginName: "plugin-1",
		callbacks:  make(map[string]func(message *runtime.MessageRequest)),
		pallet:     []reactive.Object{},
		runtime:    []reactive.Object{},
	}
	cli.getPallet()
	//cli.bus()
	cli.callbacks[createObject] = cli.addObject

	var err error

	var conn *grpc.ClientConn

	for {
		conn, err = cli.connectWithRetry()
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()
		c := runtime.NewRuntimeServiceClient(conn)
		cli.serverConnection = c
		if err := cli.registerPlugin(); err != nil {
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
