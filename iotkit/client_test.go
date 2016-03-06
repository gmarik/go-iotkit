package iotkit

import (
	"encoding/json"
	"time"

	"testing"
)

func TestNewClient_defaultClient(t *testing.T) {
	c := NewClient(nil)

	if c.client == nil {
		t.Fatal("Default Client must be created!")
	}

	if c.BaseURL.String() != "https://dashboard.us.enableiot.com/v1/api/" {
		t.Error("Invalid base url", c.BaseURL.String())
	}
}

func TestTimeJsonSerialization(t *testing.T) {

	v := struct{ Time }{Time(time.Unix(1455839028, 0))}
	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	exp := "1455839028000"
	got := string(bs)
	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
