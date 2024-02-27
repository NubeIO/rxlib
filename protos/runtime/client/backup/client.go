package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers"
	runtimeClient "github.com/NubeIO/rxlib/protos/runtime/runtimeclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// Protocol defines the interface for different protocols (gRPC, HTTP, etc.)
type Protocol interface {
	ObjectsDeploy(object *rxlib.Deploy, opts *Opts) (*rxlib.Deploy, error)
	Close() error
}

type Callback struct {
	UUID string
	Body string
}

type Opts struct {
}

// GRPCClient implements the Protocol for gRPC
type GRPCClient struct {
	client runtimeClient.RuntimeServiceClient
	conn   *grpc.ClientConn
}

func (g *GRPCClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts) (*rxlib.Deploy, error) {
	_, err := g.objectsDeploy(ObjectDeployToProto(object))
	if err != nil {
		return nil, err
	}
	return object, err

}

type MQTTClient struct {
	client   mqtt.Client
	requests map[string]chan *MQTTPayload
}

type MQTTPayload struct {
	RequestUUID string      `json:"request_uuid"`
	Payload     interface{} `json:"payload"`
}

func NewMQTTClient(client mqtt.Client) *MQTTClient {
	return &MQTTClient{
		client:   client,
		requests: make(map[string]chan *MQTTPayload),
	}
}

func (m *MQTTClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts) (*rxlib.Deploy, error) {
	// Define the request and response topics
	requestTopic := "objects/deploy/request"
	responseTopic := "objects/deploy/response"
	// Call mqttRequestResponse to send the request and wait for a response
	responsePayload, err := m.mqttRequestResponse(requestTopic, responseTopic, object, 10*time.Second)
	if err != nil {
		return nil, err
	}
	// Assert the type of responsePayload.Payload to []byte
	payload, ok := responsePayload.Payload.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid payload type")
	}
	// Unmarshal the response payload into a rxlib.Deploy object
	var response rxlib.Deploy
	if err := json.Unmarshal(payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MQTTClient) mqttRequestResponse(requestTopic, responseTopic string, payload any, timeout time.Duration) (*MQTTPayload, error) {
	// Generate a new UUID for the request
	requestUUID := helpers.UUID()

	// Create a channel to receive the response
	responseChan := make(chan *MQTTPayload, 1)
	m.requests[requestUUID] = responseChan
	defer delete(m.requests, requestUUID)

	// Subscribe to the response topic
	token := m.client.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var response MQTTPayload
		if err := json.Unmarshal(msg.Payload(), &response); err != nil {
			// Handle JSON unmarshal error
			return
		}
		if response.RequestUUID == requestUUID {
			// Send the response to the channel
			responseChan <- &response
		}
	})
	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}
	defer m.client.Unsubscribe(responseTopic)

	// Marshal the payload to JSON
	marshaledPayload, err := json.Marshal(MQTTPayload{
		RequestUUID: requestUUID,
		Payload:     payload,
	})
	if err != nil {
		return nil, err
	}

	// Publish the request
	m.client.Publish(requestTopic, 0, false, marshaledPayload)

	// Wait for the response or timeout
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("request timed out")
	}
}

func (m *MQTTClient) mqttResponseCallback(msg mqtt.Message) {
	var payload MQTTPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		fmt.Printf("Error unmarshaling MQTT payload: %v\n", err)
		return
	}

	if responseChan, ok := m.requests[payload.RequestUUID]; ok {
		responseChan <- &payload
	}
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

func (g *GRPCClient) Close() error {
	return g.conn.Close()
}

// HTTPClient implements the Protocol for HTTP (using resty)
type HTTPClient struct {
	client  *resty.Client
	baseURL string
}

func (h *HTTPClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts) (*rxlib.Deploy, error) {
	_, err := h.objectsDeploy(ObjectDeployToProto(object))
	if err != nil {
		return nil, err
	}
	return object, err
}

func (h *HTTPClient) objectsDeploy(object *runtimeClient.ObjectDeployRequest) (*runtimeClient.ObjectDeploy, error) {
	// convert the object to JSON
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request object: %v", err)
	}

	// Set up the request
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

	// Unmarshal the response into a runtimeClient.ObjectDeploy struct
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
	default:
		return nil, errors.New("unsupported protocol")
	}
}

func (m *Client) Close() {
	m.protocol.Close()
}

func (m *Client) ObjectsDeploy(object *rxlib.Deploy, opts *Opts) (*rxlib.Deploy, error) {
	return m.protocol.ObjectsDeploy(object, opts)
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
