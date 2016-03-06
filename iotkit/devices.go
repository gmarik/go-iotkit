package iotkit

import (
	"fmt"
	"net/http"
)

// https://github.com/enableiot/iotkit-api/wiki/Device-Management

type DevicesAPI struct {
	client *Client
}

// Device represents and IoT device
type Device struct {
	// ID is a unique user-provided uuid
	ID string
	// Token is a device's JWT token used to authenticate device-related API calls
	// valid for 24h
	Token string
}

type Component struct {
	ID   string `json:"cid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// CreateComponent creates component provided
// - component info
// - device's ID and auth token
// - and account's ID
// comonent was created if no error returned
// https://github.com/enableiot/iotkit-api/wiki/Device-Management#add-a-component-to-a-device
func (a *DevicesAPI) CreateComponent(comp Component, d Device, account Account) (*http.Response, error) {
	path := fmt.Sprintf("accounts/%s/devices/%s/components", account.ID, d.ID)
	req, err := a.client.NewRequest("POST", path, &comp)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.Token)

	resp, err := a.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	if err := checkStatusCreated(resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// Activate
// https://github.com/enableiot/iotkit-api/wiki/Device-Management#activate-one-device

func (dapi *DevicesAPI) ActivateDevice(activationCode string, a Account, d Device) (string, *http.Response, error) {
	path := fmt.Sprintf("accounts/%s/devices/%s/activation", a.ID, d.ID)

	body := struct {
		ActivationCode string `json:"activationCode"`
	}{activationCode}

	req, err := dapi.client.NewRequest("PUT", path, &body)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.Token)

	var ebody struct {
		DeviceToken string `json:"deviceToken"`
	}

	resp, err := dapi.client.Do(req, &ebody)
	if err != nil {
		return "", resp, err
	}

	return ebody.DeviceToken, resp, nil
}
