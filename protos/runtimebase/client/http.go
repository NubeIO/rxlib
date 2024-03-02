package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
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
	newUUID := uuid.New().String()
	go func() {
		resp, err := h.httpRequestWithTimeout("GET", endpoint, nil, opts)
		if err != nil {
			callback(newUUID, nil, err)
			return
		}

		if resp.IsError() {
			callback(newUUID, nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode()))
			return
		}
		callback(newUUID, &Message{"OK"}, nil)
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

func (h *HTTPClient) Command(opts *Opts, command *rxlib.Command, callback func(string, *rxlib.CommandResponse, error)) (string, error) {
	endpoint := "/command"
	if opts == nil {
		return "", fmt.Errorf("opts body can not be empty")
	}
	newUUID := helpers.UUID()
	go func() {
		resp, err := h.httpRequestWithTimeout("POST", endpoint, command, opts)
		if err != nil {
			callback(newUUID, nil, err)
			return
		}
		if resp.IsError() {
			callback(newUUID, nil, fmt.Errorf("HTTP request failed with status code: %d err: %s", resp.StatusCode(), resp.String()))
			return
		}
		var result *rxlib.CommandResponse
		err = json.Unmarshal(resp.Body(), &result)
		callback(newUUID, result, nil)
	}()
	return newUUID, nil

}

func (h *HTTPClient) makeRequestWithCallback(method, endpoint string, body interface{}, opts *Opts, callback func(*Response, error)) (string, error) {
	newUUID := uuid.New().String() // Generate a UUID for the request

	return newUUID, nil
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

func (h *HTTPClient) objectsDeploy(object *runtime.ObjectConfig) (*runtime.ObjectConfig, error) {
	// convert the object to JSON
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request object: %v", err)
	}

	// set up the request
	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonObject).
		Post(fmt.Sprintf("/%s", "runtimebase/deploy"))

	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode())
	}

	// unmarshal the response into a runtimeClient.ObjectDeploy struct
	var result runtime.ObjectConfig
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
