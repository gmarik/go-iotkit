package iotkit

import (
	"fmt"
	"net/http"
)

/// https://github.com/enableiot/iotkit-api/wiki/Component-Types-Catalog

type ComponentTypesAPI struct {
	client *Client
}

// https://github.com/enableiot/iotkit-api/wiki/Component-Types-Catalog#create-a-new-custom-component-type
type ComponentType struct {
	ID string `json:"id,omitempty"`

	// Component Type Name  e.g., "temperature"
	Dimension string `json:"dimension"`

	// API version  1.0 | 2.0
	Version string `json:"version"`
	// Component Type  sensor | actuator
	Type string `json:"type"`

	// Data Type:  Number | String | Boolean | ByteArray
	DataType string `json:"dataType"`
	// Data Format:  float | boolean | string | percentage | integer
	Format string `json:"format"`

	// Minimum Value  e.g., -150
	// Maximum Value  e.g., 150
	Min float64 `json:"min"`
	Max float64 `json:"max"`

	// Units Name  e.g., "Degrees Celsius"
	MeasureUnit string `json:"measureunit"`

	// Data Series Type
	// "timeSeries" for Numbers
	// "rawData" for String and Boolean
	// "binaryDataRenderer" for ByteString
	Display string `json:"display"`

	// component URI
	HRef string `json:"href,omitempty"`
}

// CreateComponentType creates custom component type
// returns component with the ID assigned
func (a *ComponentTypesAPI) CreateComponentType(account Account, ct ComponentType) (*ComponentType, *http.Response, error) {
	path := fmt.Sprintf("/accounts/%s/cmpcatalog", account.ID)
	req, err := a.client.NewRequest("POST", path, &ct)
	if err != nil {
		return &ct, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+account.Token)

	resp, err := a.client.Do(req, &ct)
	if err != nil {
		return &ct, resp, err
	}

	if err := checkStatusCreated(resp); err != nil {
		return &ct, resp, err
	}

	return &ct, resp, nil
}
