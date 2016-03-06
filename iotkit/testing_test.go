package iotkit

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url

}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func dumpRequest(resp *http.Request, t *testing.T) {
	dump, err := httputil.DumpRequest(resp, true)
	if err != nil {
		t.Fatal("Error Dumping Request: ", err)
	}

	t.Log("\n", string(dump))
}
func dumpResponse(resp *http.Response, t *testing.T) {
	dump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		t.Fatal("Error Dumping Response: ", err)
	}

	t.Log("\n", string(dump))
}

func assertHeader(t *testing.T, key, value string, headers http.Header) {
	exp := value
	got := headers.Get(key)
	if got != exp {
		t.Errorf("AssertHeader:\n\nExp: %v\nGot: %v", exp, got)
	}
}
