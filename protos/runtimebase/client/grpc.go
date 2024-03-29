package client

import (
	"context"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"
)

type GRPCClient struct {
	client runtime.RuntimeServiceClient
	conn   *grpc.ClientConn
}

func (g *GRPCClient) Command(opts *Opts, command *rxlib.ExtendedCommand, callback func(string, *runtime.CommandResponse, error)) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GRPCClient) command(object *runtime.Command) (*runtime.Command, error) {
	return nil, nil
}

func (g *GRPCClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	uuid := uuid.New().String()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		resp, err := g.client.Ping(ctx, &runtime.PingRequest{})
		var message *Message
		if err == nil && resp != nil {
			message = &Message{Message: resp.GetMessage()}
		}
		callback(uuid, message, err)
	}()
	return uuid, nil
}

func (g *GRPCClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	uuid := uuid.New().String()
	go func() {
		//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		//defer cancel()
		//
		//resp, err := g.client.ObjectsDeploy(ctx, ObjectDeployToProto(object))
		//callback(&Callback{UUID: uuid, Body: resp}, err)
	}()
	return uuid, nil
}

func (g *GRPCClient) Close() error {
	return g.conn.Close()
}

func (g *GRPCClient) objectsDeploy(object *runtime.ObjectDeploy) (*runtime.ObjectDeploy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	r, err := g.client.ObjectsDeploy(ctx, object)
	if err != nil {
		return nil, fmt.Errorf("could not deploy objects: %v", err)
	}
	return r, nil
}
