package iotkit

import (
	"net/http"
	// "reflect"
	"reflect"
	"testing"
	"time"

	"io/ioutil"

	// "fmt"
)

func TestObservations_Create_Success(t *testing.T) {
	setup()
	defer teardown()

	var (
		device = Device{
			ID:    "an_id",
			Token: "deviceToken",
		}

		gotReqBody []byte
		expReqBody = `{"accountId":"11111111-1111-1111-1111-111111111111","on":1455743045000,"data":[{"componentId":"cccccccc-cccc-cccc-cccc-cccccccccccc","on":1455743045000,"value":"21.7","attributes":{"i":0}}]}`

		batch = ObservationBatch{
			AccountId: "11111111-1111-1111-1111-111111111111",
			On:        Time(time.Unix(1455743045000/1000, 0)),
			Data: []Observation{
				//{"on": 1455743045000, "componentId": "cccccccc-cccc-cccc-cccc-cccccccccccc", "value": "21.7", "attributes": {"i": 0}}]}
				{
					ComponentId: "cccccccc-cccc-cccc-cccc-cccccccccccc",
					On:          Time(time.Unix(1455743045000/1000, 0)),
					Value:       "21.7",
					Attributes: map[string]interface{}{
						"i": 0,
					},
				},
			},
		}
	)

	mux.HandleFunc("/data/"+device.ID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		gotReqBody, _ = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
	})

	_, err := client.ObservationAPI.Create(batch, device)

	{
		exp := append([]byte(expReqBody), 10)
		got := gotReqBody
		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}

	if err != nil {
		t.Error(err)
	}
}
