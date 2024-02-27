package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/rxlib"
	runtimeClient "github.com/NubeIO/rxlib/protos/runtime/runtimeclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"time"
)

type Protocol interface {
	ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error)
	Close() error
}

type Callback struct {
	UUID string
	Body interface{}
}

type Opts struct {
	// Your options
}

type GRPCClient struct {
	client runtimeClient.RuntimeServiceClient
	conn   *grpc.ClientConn
}

func (g *GRPCClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	uuid := uuid.New().String()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		resp, err := g.client.ObjectsDeploy(ctx, ObjectDeployToProto(object))
		callback(&Callback{UUID: uuid, Body: resp}, err)
	}()
	return uuid, nil
}

func (g *GRPCClient) Close() error {
	return g.conn.Close()
}

func (g *GRPCClient) objectsDeploy(object *runtimeClient.ObjectDeployRequest) (*runtimeClient.ObjectDeploy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration(object.Timeout))
	defer cancel()
	r, err := g.client.ObjectsDeploy(ctx, object)
	if err != nil {
		return nil, fmt.Errorf("could not deploy objects: %v", err)
	}
	return r.ObjectDeploy, nil
}

type MQTTClient struct {
	mqtt.Client
	requests map[string]chan *MQTTPayload
}

func (m *MQTTClient) Close() error {
	//TODO implement me
	panic("implement me")
}

type MQTTPayload struct {
	RequestUUID string      `json:"request_uuid"`
	Payload     interface{} `json:"payload"`
}

func newMQTTClient(opts *mqtt.ClientOptions) *MQTTClient {
	return &MQTTClient{
		Client:   mqtt.NewClient(opts),
		requests: make(map[string]chan *MQTTPayload),
	}
}

func (m *MQTTClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	// Define the request and response topics
	requestTopic := "objects/deploy/request"
	responseTopic := "objects/deploy/response"

	uuid := uuid.New().String()

	// Create a channel to receive the response
	responseChan := make(chan *MQTTPayload, 1)
	m.requests[uuid] = responseChan
	defer delete(m.requests, uuid)

	// Subscribe to the response topic
	token := m.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var response MQTTPayload
		if err := json.Unmarshal(msg.Payload(), &response); err != nil {
			// Handle JSON unmarshal error
			return
		}
		if response.RequestUUID == uuid {
			// Send the response to the channel
			responseChan <- &response
		}
	})
	token.Wait()
	if token.Error() != nil {
		return "", token.Error()
	}
	defer m.Unsubscribe(responseTopic)

	// Marshal the payload to JSON
	marshaledPayload, err := json.Marshal(MQTTPayload{
		RequestUUID: uuid,
		Payload:     object,
	})
	if err != nil {
		return "", err
	}

	// Publish the request
	m.Publish(requestTopic, 0, false, marshaledPayload)

	// Start a goroutine to handle the response
	go func() {
		select {
		case response := <-responseChan:
			// Unmarshal the response payload into a rxlib.Deploy object
			var responseObject rxlib.Deploy
			if err := json.Unmarshal(response.Payload.([]byte), &responseObject); err != nil {
				callback(nil, err)
				return
			}
			callback(&Callback{UUID: uuid, Body: &responseObject}, nil)
		case <-time.After(10 * time.Second): // Adjust the timeout as needed
			callback(nil, fmt.Errorf("request timed out"))
		}
	}()

	return uuid, nil
}

// HTTPClient implements the Protocol for HTTP (using resty)
type HTTPClient struct {
	client  *resty.Client
	baseURL string
}

func (h *HTTPClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	uuid := uuid.New().String()
	go func() {
		// Assuming you have a method `httpObjectsDeploy` that makes the HTTP request
		resp, err := h.objectsDeploy(ObjectDeployToProto(object))
		callback(&Callback{UUID: uuid, Body: resp}, err)
	}()
	return uuid, nil
}

func (h *HTTPClient) objectsDeploy(object *runtimeClient.ObjectDeployRequest) (*runtimeClient.ObjectDeploy, error) {
	// convert the object to JSON
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request object: %v", err)
	}

	// set up the request
	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonObject).
		Post(fmt.Sprintf("/%s/%s", h.baseURL, "runtime/deploy"))

	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode())
	}

	// unmarshal the response into a runtimeClient.ObjectDeploy struct
	var result runtimeClient.ObjectDeploy
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

func (h *HTTPClient) Close() error {
	// Implement any necessary cleanup for the HTTP client
	return nil
}

// Client struct now holds a Protocol instead of the concrete client
type Client struct {
	protocol Protocol
}

func NewClient(protocol string, port, httpPort int) (*Client, error) {
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
		client.BaseURL = fmt.Sprintf("http:localhost:%d/api", httpPort)
		return &Client{protocol: &HTTPClient{client: client}}, nil
	case "mqtt":

		return &Client{protocol: &MQTTClient{}}, nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}

// Add other missing methods similar to the above implementations
func (m *Client) Close() {
	m.protocol.Close()
}

func (m *Client) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	return m.protocol.ObjectsDeploy(object, opts, callback)
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
