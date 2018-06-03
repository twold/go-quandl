package request

import (
	"io"
	"io/ioutil"
	"net/http"
)

// A Request is the service request to be made.
type Request struct {
	HTTPRequest  *http.Request
	HTTPResponse *http.Response
	Body         io.ReadCloser
	Params       interface{}
	Error        error
	Data         interface{}
}

func New(method, Url string, body io.Reader) *Request {
	// initialize new Request
	request := &Request{}
	if method == "" {
		method = "GET"
	}

	// build http request and save to output
	request.HTTPRequest, request.Error = http.NewRequest(method, Url, body)
	if request.Error != nil {
		return request
	}

	request.HTTPRequest.Header.Add("cache-control", "no-cache")

	request.HTTPResponse, request.Error = http.DefaultClient.Do(request.HTTPRequest)
	if request.Error != nil {
		return request
	}
	// save response body as readcloser
	request.Body = ioutil.NopCloser(request.HTTPResponse.Body)
	return request
}
