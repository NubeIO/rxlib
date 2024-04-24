package restc

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Options struct {
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type Rest interface {
	Execute(method, endPoint string, opts *Options) *Response
}

func New() Rest {
	return &impl{client: resty.New()}
}

type impl struct {
	client *resty.Client
}

func (i *impl) Execute(method, endPoint string, opts *Options) *Response {
	req := i.client.R()
	if opts != nil {
		// Set headers if provided
		if opts.Headers != nil {
			req.SetHeaders(opts.Headers)
		}

		// Set body if provided
		if opts.Body != nil {
			req.SetBody(opts.Body)
		}
	}

	// Execute request
	resp, err := req.Execute(method, endPoint)

	// Create response struct
	return &Response{
		RawResponse: resp,
		Error:       err,
	}
}

type Response struct {
	RawResponse *resty.Response
	Error       error
}

func (r *Response) Code() int {
	if r.RawResponse != nil {
		return r.RawResponse.StatusCode()
	}
	return 0
}

func (r *Response) GetError() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return ""
}

func (r *Response) IsError() bool {
	return r.Error != nil
}

func (r *Response) AsJSON(target interface{}) error {
	if r.RawResponse == nil {
		return r.Error
	}
	return json.Unmarshal(r.RawResponse.Body(), target)
}

func (r *Response) AsString() string {
	if r.RawResponse != nil {
		return r.RawResponse.String()
	}
	return ""
}
