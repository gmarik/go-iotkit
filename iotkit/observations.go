package iotkit

import "net/http"

type ObservationAPI struct {
	client *Client
}

type ObservationBatch struct {
	AccountId string `json:"accountId"`
	On        Time   `json:"on"`
	// Loc       [2]float32    `json:"loc,omitempty"`
	Data []Observation `json:"data"`
}

type Observation struct {
	ComponentId string                 `json:"componentId"`
	On          Time                   `json:"on"`
	Value       string                 `json:"value"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

func (c *ObservationAPI) Create(data ObservationBatch, d Device) (*http.Response, error) {

	req, err := c.client.NewRequest("POST", "data/"+d.ID, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+d.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	if err := checkStatusCreated(resp); err != nil {
		return resp, err
	}

	return resp, nil
}
