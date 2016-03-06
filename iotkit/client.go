package iotkit

import (
	"encoding/json"

	"bytes"

	"net/http"
	"net/url"
	// "net/http/httputil"

	"fmt"
	"io"
	"io/ioutil"
)

const (
	ApiHost = "dashboard.us.enableiot.com"
	ApiRoot = "/v1/api/"
)

type ErrorResponse struct {
	Code    int
	Status  int
	Message string
	Errors  []string
}

func (m *ErrorResponse) Error() string {
	return fmt.Sprintf("Code: %d; Status: %d; Message: %q; Errors: %v", m.Code, m.Status, m.Message, m.Errors)
}

func ApiPath(path string) string {
	return "https://" + ApiHost + ApiRoot + path
}

// https://0value.com/Let-the-Doer-Do-it
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	client  Doer
	BaseURL *url.URL

	*AuthorizationApi
	*ObservationAPI
	*DevicesAPI
	*ComponentTypesAPI
}

func NewClient(client Doer) *Client {

	if client == nil {
		client = &http.Client{}
	}

	//always correct since they're predefined constants
	u, _ := url.Parse("https://" + ApiHost + ApiRoot)

	c := &Client{
		client:  client,
		BaseURL: u,
	}

	c.AuthorizationApi = &AuthorizationApi{client: c}
	c.ObservationAPI = &ObservationAPI{client: c}
	c.DevicesAPI = &DevicesAPI{client: c}
	c.ComponentTypesAPI = &ComponentTypesAPI{client: c}

	return c
}

func (c *Client) NewRequest(method, relUrl string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(relUrl)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		// fmt.Printf("Encoding: %#v\n", body)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}
	return resp, err
}

func checkResponse(r *http.Response) error {
	var c = r.StatusCode
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if http.StatusUnauthorized == c {
		return fmt.Errorf("Error: %s", body)
	}

	errorResponse := &ErrorResponse{}
	if err := json.Unmarshal(body, errorResponse); err != nil {
		return err
	}

	return errorResponse
}

func checkStatusCreated(resp *http.Response) error {
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Expected Status: %q; Got: %q", "201 Created", resp.Status)
	}
	return nil
}
