package client

import (
	"errors"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib"
	runtimeClient "github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Protocol interface {
	ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error)
	Close() error
	Ping(opts *Opts, callback func(string, *Message, error)) (string, error)
	Command(opts *Opts, command *rxlib.ExtendedCommand, callback func(string, *rxlib.CommandResponse, error)) (string, error)
}

const defaultTimeout = 2

type Callback struct {
	UUID string
	Body interface{}
}

type Opts struct {
	RequestUUID    string
	TargetGlobalID string
	SenderGlobalID string
	Timeout        time.Duration
	Headers        map[string]string
}

type Message struct {
	Message string
}

// Client struct now holds a Protocol instead of the concrete client
type Client struct {
	protocol Protocol
}

func (m *Client) Command(opts *Opts, command *rxlib.ExtendedCommand, callback func(string, *rxlib.CommandResponse, error)) (string, error) {
	return m.protocol.Command(opts, command, callback)
}

func (m *Client) Close() error {
	//TODO implement me
	panic("implement me")
}

func NewClient(ip, protocol string, port, httpPort int, mqtt mqttwrapper.MQTT) (Protocol, error) {
	if ip == "" {
		ip = "localhost"
	}
	switch protocol {
	case "grpc":
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("did not connect: %v", err)
		}
		c := runtimeClient.NewRuntimeServiceClient(conn)
		return &Client{protocol: &GRPCClient{client: c, conn: conn}}, nil
	case "http":
		client := resty.New()
		client.BaseURL = fmt.Sprintf("http://%s:%d/api", ip, httpPort)
		return &Client{protocol: &HTTPClient{client: client}}, nil
	case "mqtt":
		c, _ := newMQTTClient(mqtt)
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

//func ObjectDeployToProto(obj *rxlib.Deploy) *runtimeClient.ObjectDeploy {
//	a := &runtimeClient.ObjectDeploy{
//		Deleted: obj.Deleted,
//		New:     ObjectsConfigToProto(obj.New),
//		Updated: ObjectsConfigToProto(obj.Updated),
//	}
//	return a
//}

//func ObjectsConfigToProto(objs []*runtime.ObjectConfig) []*runtimeClient.ObjectConfig {
//	var out []*runtimeClient.ObjectConfig
//	for _, obj := range objs {
//		out = append(out, ObjectConfigToProto(obj))
//	}
//	out = append(out)
//	return out
//}

//func ObjectConfigToProto(obj *runtime.ObjectConfig) *runtimeClient.ObjectConfig {
//	return &runtimeClient.ObjectConfig{
//		Id:          obj.ID,
//		Info:        nil,
//		Inputs:      nil,
//		Outputs:     nil,
//		Meta:        nil,
//		Stats:       nil,
//		Connections: nil,
//	}
//}
