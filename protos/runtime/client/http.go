package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/protos/runtime/protoruntime"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"time"
)

// HTTPClient implements the Protocol for HTTP (using resty)
type HTTPClient struct {
	client  *resty.Client
	baseURL string
}

type Response struct {
	UUID string
	Body interface{}
}

func (h *HTTPClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	endpoint := "/ping"

	return h.makeRequestWithCallback("GET", endpoint, nil, opts, func(response *Response, err error) {
		var message *Message
		if err == nil && response != nil {
			err = json.Unmarshal(response.Body.([]byte), &message)
		}
		callback(response.UUID, message, err)
	})
}

func (h *HTTPClient) makeRequestWithCallback(method, endpoint string, body interface{}, opts *Opts, callback func(*Response, error)) (string, error) {
	newUUID := uuid.New().String() // Generate a UUID for the request

	go func() {
		resp, err := h.httpRequestWithTimeout(method, endpoint, body, opts)
		if err != nil {
			callback(nil, err)
			return
		}

		if resp.IsError() {
			callback(nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode()))
			return
		}

		callback(&Response{UUID: newUUID, Body: resp.Body()}, nil)
	}()

	return newUUID, nil
}

func (h *HTTPClient) httpRequestWithTimeout(method, endpoint string, body interface{}, opts *Opts) (*resty.Response, error) {
	request := h.client.R()
	if body != nil {
		request.SetBody(body)
	}
	var timeout time.Duration
	if opts != nil && opts.Timeout != 0 {
		timeout = opts.Timeout
	} else {
		timeout = 2 * time.Second // Default timeout
	}

	if opts != nil && opts.Headers != nil {
		for key, value := range opts.Headers {
			request.SetHeader(key, value)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request.SetContext(ctx)

	var resp *resty.Response
	var err error
	switch method {
	case "GET":
		resp, err = request.Get(endpoint)
	case "POST":
		resp, err = request.Post(endpoint)
	case "PATCH":
		resp, err = request.Patch(endpoint)
	case "PUT":
		resp, err = request.Put(endpoint)
	case "DELETE":
		resp, err = request.Delete(endpoint)
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

func (h *HTTPClient) objectsDeploy(object *protoruntime.ObjectDeployRequest) (*protoruntime.ObjectDeploy, error) {
	// convert the object to JSON
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request object: %v", err)
	}

	// set up the request
	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonObject).
		Post(fmt.Sprintf("/%s", "runtime/deploy"))

	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode())
	}

	// unmarshal the response into a runtimeClient.ObjectDeploy struct
	var result protoruntime.ObjectDeploy
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
