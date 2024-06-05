package extensionlib

import (
	"context"
	"fmt"
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

func (inst *Extensions) Register() error {
	var err error

	var conn *grpc.ClientConn
	go inst.server.Run()
	for {
		conn, err = inst.connectWithRetry()
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer conn.Close()
		messages[helpers.UUID()] = fmt.Sprintf("connected to server")
		c := runtime.NewRuntimeServiceClient(conn)
		inst.grpcClient = c
		if err := inst.registerExtension(); err != nil {
			log.Fatalf("could not register plugin: %v", err)
		}
		messages[helpers.UUID()] = fmt.Sprintf("registered ExtensionPayload")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := inst.startServerStreaming(ctx, conn); err != nil {
			log.Printf("streaming error: %v, reconnecting...", err)
			continue
		}
	}
}

// registerExtension register the Extensions with the server
func (inst *Extensions) registerExtension() error {
	s := inst.grpcClient
	regCtx, regCancel := context.WithTimeout(context.Background(), time.Second)
	defer regCancel()
	info := &runtime.Extension{Name: "ExampleExtension", Uuid: inst.name, Pallet: reactive.ConvertObjects(inst.pallet)}
	_, err := s.RegisterExtension(regCtx, info)
	if err != nil {
		return fmt.Errorf("could not register ExtensionPayload: %v", err)
	}
	fmt.Println("Registered ExtensionPayload")
	return nil
}

func (inst *Extensions) sendMessageToServer(content, key string) error {
	if inst.stream == nil {
		messages[helpers.UUID()] = fmt.Sprintf("stream is not initialized %s %s", content, key)
		return fmt.Errorf("stream is not initialized")
	}
	messages[helpers.UUID()] = fmt.Sprintf("new stream message %s %s", content, key)
	if err := inst.stream.Send(&runtime.MessageRequest{ExtensionUUID: inst.name, Key: key, StringPayload: content}); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// newObject message will come from the server to create an instance of the object
func (inst *Extensions) newObject(message *runtime.MessageRequest) {
	messages[helpers.UUID()] = "Add new object"
	object := message.GetObject()
	if object == nil {
		messages[helpers.UUID()] = "Add new object, object was empty"
		if err := inst.sendMessageToServer("failed to get object", "error"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
		return
	}
	if object.GetInfo() == nil {
		messages[helpers.UUID()] = "object info were empty"
		return
	}
	objectID := object.GetInfo().GetObjectID()
	if objectID == "" {
		if err := inst.sendMessageToServer("objectID is empty", "error"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
		messages[helpers.UUID()] = "object id was empty"
		return
	}
	var objectExists bool
	for _, obj := range inst.pallet {
		if objectID == obj.GetInfo().GetObjectID() {
			objectExists = true
		}
	}

	if !objectExists {
		messages[helpers.UUID()] = "Add new object, object was not found in pallet"
		if err := inst.sendMessageToServer("failed to find object in pallet", "error"); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
	}

	//instance := inst.objectInstance(reactive.ConvertObjectConfig(object), inst.outputCallback)
	//
	//if instance != nil {
	//	inst.runtime = append(inst.runtime, instance)
	//	inst.callbacks[instance.GetMeta().GetObjectUUID()] = instance.Handler
	//	if err := inst.sendMessageToServer("response for 1", "error"); err != nil {
	//		fmt.Printf("failed to send message: %v\n", err)
	//	}
	//	fmt.Printf("objects count: %d\n", len(inst.runtime))
	//	return
	//} else {
	//	if err := inst.sendMessageToServer("failed to find object in pallet", "error"); err != nil {
	//		fmt.Printf("failed to send message: %v\n", err)
	//	}
	//}
}

// outputCallback send a message back to the server when the output value of the object is updated
func (inst *Extensions) outputCallback(cmd *runtime.Command) {
	if err := inst.stream.Send(&runtime.MessageRequest{
		Key:           "invoke",
		ExtensionUUID: inst.name,
		Command:       cmd,
	}); err != nil {
		messages[helpers.UUID()] = fmt.Sprintf("outputCallback err: %v", err)
	}

}

// objectInstance create a new instance
//func (inst *Extensions) objectInstance(obj *reactive.BaseObject, outputUpdated func(message *runtime.Command)) reactive.Object {
//	instance := subtract.New(outputUpdated)
//	base := instance.New(obj)
//	return base
//}

var infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

// startServerStreaming stream messages from the server
func (inst *Extensions) startServerStreaming(ctx context.Context, conn *grpc.ClientConn) error {
	c := runtime.NewRuntimeServiceClient(conn)

	// Start bidirectional streaming with the given context
	stream, err := c.ExtensionStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to start streaming: %v", err)
	}
	inst.stream = stream
	// Send a message to the server. initiate the client connection to the server. The server will persist the client
	if err := stream.Send(&runtime.MessageRequest{ExtensionUUID: inst.name, StringPayload: "brith from ExtensionPayload"}); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	messages["start"] = time.Now()
	// Receive messages from the server
	for {
		select {
		case <-ctx.Done():
			messages[helpers.UUID()] = fmt.Sprintf("exit stream loop")
			return nil // Context canceled, exit the loop
		default:
			in, err := stream.Recv()
			messages[helpers.UUID()] = fmt.Sprintf("new stream message key: %s", in.Key)
			if err == io.EOF {
				return nil // Server closed the stream
			}
			if err != nil {
				return fmt.Errorf("failed to receive message: %v", err)
			}
			if in.Key == "create-object" {
				inst.newObject(in)
			}
			if in.Key == "input-updated" {
				if callback, ok := inst.callbacks[in.GetObjectUUID()]; ok {
					callback(in)
				} else {
					fmt.Printf("Received message from server unknown key: %s\n", in.Key)
				}
			}

		}
	}
}

func (inst *Extensions) connectWithRetry() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	var count int
	for {
		target := fmt.Sprintf("localhost:%s", inst.grpcPort)
		conn, err = grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("could not connect: %v", err)
			messages[helpers.UUID()] = fmt.Sprintf("connectWithRetry count: %d", count)
			time.Sleep(30 * time.Second) // Retry after 30 seconds
			continue
		}
		break
	}

	return conn, nil
}
