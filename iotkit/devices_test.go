package iotkit

import (
	"net/http"
	"reflect"
	"testing"

	"fmt"
)

var (
	deviceId    = "6F4DBDD2-0C89-4B8B-BB8A-C3CD8A71B93A"
	accountId   = "E6A8EC94-2C25-4A30-8C07-CCCC4E227F06"
	componentId = "BE02BAB0-0D39-4988-BF51-0A520A4E31EA"
)

func TestDevices_CreateComponent(t *testing.T) {

	tcases := []struct {
		h   http.HandlerFunc
		err error
	}{
		{
			err: &ErrorResponse{Code: 400, Status: 400, Message: "Invalid request"},
			h: func(w http.ResponseWriter, r *http.Request) {
				assertHeader(t, "Authorization", "Bearer ASECRETTOKEN", r.Header)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"code":400,"status":400,"message":"Invalid request"}`)
			},
		},

		{
			err: nil,
			h: func(w http.ResponseWriter, r *http.Request) {
				assertHeader(t, "Authorization", "Bearer ASECRETTOKEN", r.Header)
				w.WriteHeader(http.StatusCreated)
			},
		},
	}

	c := Component{
		ID:   componentId,
		Name: "test component",
		Type: "custom.v1.0",
	}

	d := Device{
		ID:    deviceId,
		Token: "ASECRETTOKEN",
	}

	a := Account{
		ID: accountId,
	}

	for _, tcase := range tcases {
		// TODO: make stateless
		setup()
		defer teardown()

		path := fmt.Sprintf("/accounts/%s/devices/%s/components", accountId, deviceId)
		mux.HandleFunc(path, tcase.h)

		_, err := client.CreateComponent(c, d, a)

		got, exp := err, tcase.err

		if !reflect.DeepEqual(err, exp) {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
}
