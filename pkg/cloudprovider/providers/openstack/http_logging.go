package openstack

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/rackspace/gophercloud"
)

type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(resp *http.Response, err error)
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
	glog.V(2).Infof("---> %s %s", req.Method, req.URL)
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(resp *http.Response, err error) {
	glog.V(2).Infof("<--- %d %s err=%v, header=%v, body=%v", resp.StatusCode, resp.Request.URL, err, resp.Header, resp.Body)
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		t.logResponse(resp, err)
		return resp, err
	}

	t.logResponse(resp, nil)

	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(resp *http.Response, err error) {
	if t.LogResponse != nil {
		t.LogResponse(resp, err)
	} else {
		DefaultLogResponse(resp, err)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}

func EnableHTTPLogging(client *gophercloud.ServiceClient) {
	httpClient := http.Client{
		Transport: &Transport{},
	}
	client.HTTPClient = httpClient
}
