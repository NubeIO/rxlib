package client

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rxlib"
	runtimeClient "github.com/NubeIO/rxlib/protos/runtime/protoruntime"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Protocol interface {
	ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error)
	Close() error
	Ping(opts *Opts, callback func(string, *Message, error)) (string, error)
}

const defaultTimeout = 2

type Callback struct {
	UUID string
	Body interface{}
}

type Opts struct {
	Timeout time.Duration
	Headers map[string]string
}

type Message struct {
	Message string
}

// Client struct now holds a Protocol instead of the concrete client
type Client struct {
	protocol Protocol
}

func (m *Client) Close() error {
	//TODO implement me
	panic("implement me")
}

func NewClient(protocol string, port, httpPort int) (Protocol, error) {
	switch protocol {
	case "grpc":
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("did not connect: %v", err)
		}
		c := runtimeClient.NewRuntimeServiceClient(conn)
		return &Client{protocol: &GRPCClient{client: c, conn: conn}}, nil
	case "http":
		client := resty.New()
		client.BaseURL = fmt.Sprintf("http://localhost:%d/api", httpPort)
		return &Client{protocol: &HTTPClient{client: client}}, nil
	case "mqtt":
		c, _ := newMQTTClient()
		return &Client{protocol: c}, nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}

func (m *Client) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	return m.protocol.ObjectsDeploy(object, opts, callback)
}

func (m *Client) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	return m.protocol.Ping(opts, callback)
}

func duration(timeout int32) time.Duration {
	if timeout == 0 {
		timeout = 1
	}
	return time.Second * time.Duration(timeout)
}

func ObjectDeployToProto(obj *rxlib.Deploy) *runtimeClient.ObjectDeployRequest {
	a := &runtimeClient.ObjectDeploy{
		Deleted: obj.Deleted,
		New:     ObjectsConfigToProto(obj.New),
		Updated: ObjectsConfigToProto(obj.Updated),
	}
	return &runtimeClient.ObjectDeployRequest{
		ObjectDeploy: a,
		Timeout:      0,
	}
}

func ObjectsConfigToProto(objs []*rxlib.ObjectConfig) []*runtimeClient.Object {
	var out []*runtimeClient.Object
	for _, obj := range objs {
		out = append(out, ObjectConfigToProto(obj))
	}
	out = append(out)
	return out
}

func ObjectConfigToProto(obj *rxlib.ObjectConfig) *runtimeClient.Object {
	return &runtimeClient.Object{
		Id:          obj.ID,
		Info:        nil,
		Inputs:      nil,
		Outputs:     nil,
		Meta:        nil,
		Stats:       nil,
		Connections: nil,
	}
}
