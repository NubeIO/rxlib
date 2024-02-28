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
	"github.com/google/uuid"
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
	UUID    string
	Message string
}

type GRPCClient struct {
	client runtimeClient.RuntimeServiceClient
	conn   *grpc.ClientConn
}

func (g *GRPCClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	//TODO implement me
	panic("implement me")
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
	RequestUUID string      `json:"requestUUID"`
	Payload     interface{} `json:"payload"`
}

func newMQTTClient() (*MQTTClient, error) {
	// Create MQTT client options
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	// Create MQTT client
	client := mqtt.NewClient(opts)
	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker: %v", token.Error())
	}

	return &MQTTClient{
		Client:   client,
		requests: make(map[string]chan *MQTTPayload),
	}, nil
}

func (m *MQTTClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	requestTopic := "ros/api/RX-1/ping"
	payloadData := &rxlib.Command{
		SenderGlobalID: "RX-1",
		Key:            "ping",
	}
	return m.RequestResponse(requestTopic, payloadData, func(uuid string, payload *Payload, err error) {
		var message *Message
		if err == nil && payload != nil {
			err = json.Unmarshal(payload.Payload, &message)
		}
		callback(uuid, message, err)
	})
}

type Payload struct {
	Payload []byte
	Topic   string
	UUID    string
}

func (m *MQTTClient) RequestResponse(requestTopic string, payloadData interface{}, callback func(string, *Payload, error)) (string, error) {
	newUUID := helpers.UUID()
	requestTopicWithUUID := fmt.Sprintf("%s_%s", requestTopic, newUUID)
	respTopicWithUUID := fmt.Sprintf("%s/response", requestTopicWithUUID)
	// Channel to signal the receipt of the message
	done := make(chan struct{})

	// Subscribe to the response topic
	token := m.Subscribe(respTopicWithUUID, 0, func(client mqtt.Client, msg mqtt.Message) {
		response := &Payload{
			Payload: msg.Payload(),
			Topic:   msg.Topic(),
		}
		_, requestUUID, err := ExtractApiTopicPath(msg.Topic())
		if err != nil {
			return
		}
		if requestUUID == newUUID {
			// Handle the response
			callback(requestUUID, response, nil)
			close(done) // Signal that the message has been received
			return
		}
	})
	token.Wait()
	if token.Error() != nil {
		return "", token.Error()
	}
	defer m.Unsubscribe(respTopicWithUUID)

	// Marshal the payload to JSON
	marshaledPayload, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}

	// Publish the request
	m.Publish(requestTopicWithUUID, 0, false, marshaledPayload)

	// Wait for response or timeout
	select {
	case <-done:
		// Message received
	case <-time.After(2 * time.Second):
		// Timeout occurred
		callback("", nil, fmt.Errorf("timeout occurred"))
	}

	return newUUID, nil
}

func (m *MQTTClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	//uuid := uuid.New().String()
	//
	//requestTopic := fmt.Sprintf("objects/deploy/%s/request", uuid)
	//responseTopic := fmt.Sprintf("objects/deploy/%s/response", uuid)
	//
	//// Subscribe to the response topic
	//token := m.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
	//	var response MQTTPayload
	//	if err := json.Unmarshal(msg.Payload(), &response); err != nil {
	//		// Handle JSON unmarshal error
	//		return
	//	}
	//	if response.RequestUUID == uuid {
	//		// Handle the response
	//		callback(response.RequestUUID, response.Payload, nil)
	//	}
	//})
	//token.Wait()
	//if token.Error() != nil {
	//	return "", token.Error()
	//}
	//defer m.Unsubscribe(responseTopic)
	//
	//// Marshal the payload to JSON
	//marshaledPayload, err := json.Marshal(MQTTPayload{
	//	RequestUUID: uuid,
	//	Payload:     object,
	//})
	//if err != nil {
	//	return "", err
	//}
	//
	//// Publish the request
	//m.Publish(requestTopic, 0, false, marshaledPayload)
	//
	//// Wait for response or timeout
	//select {
	//case <-time.After(2 * time.Second):
	//	// Timeout occurred
	//	callback("", nil, fmt.Errorf("timeout occurred"))
	//}

	return "uuid", nil
}

// HTTPClient implements the Protocol for HTTP (using resty)
type HTTPClient struct {
	client  *resty.Client
	baseURL string
}

func (h *HTTPClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	go func() {
		resp, err := h.httpRequestWithTimeout("GET", "/ping", nil, opts)
		if err != nil {
			callback("", nil, err)
			return
		}

		if resp.IsError() {
			callback("", nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode()))
			return
		}

		var result *Message
		err = json.Unmarshal(resp.Body(), &result)
		if err != nil {
			callback("", nil, fmt.Errorf("error unmarshaling response: %v", err))
			return
		}

		callback("", result, nil)
	}()

	return "", nil
}

func (h *HTTPClient) httpRequestWithTimeout(method, endpoint string, body interface{}, opts *Opts) (*resty.Response, error) {
	request := h.client.R()
	if body != nil {
		request.SetBody(body)
	}
	var timeout time.Duration
	if opts != nil {
		if opts.Timeout == 0 {
			opts.Timeout = duration(2)
		}

		if opts.Headers != nil {
			for key, value := range opts.Headers {
				request.SetHeader(key, value)
			}
		}
	} else {
		timeout = duration(defaultTimeout)
	}

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Make sure to cancel the context to release resources

	// Use the context in the request
	request.SetContext(ctx)

	var resp *resty.Response
	var err error
	switch method {
	case "GET":
		resp, err = request.Get(h.baseURL + endpoint)
	case "POST":
		resp, err = request.Post(h.baseURL + endpoint)
	case "PUT":
		resp, err = request.Put(h.baseURL + endpoint)
	case "DELETE":
		resp, err = request.Delete(h.baseURL + endpoint)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	return resp, err
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

//// Close Add other missing methods similar to the above implementations
//func (m *Client) Close() {
//	m.protocol.Close()
//}

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
