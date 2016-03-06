package iotkit

import (
	"net/http"
	"reflect"
	"testing"

	"fmt"
)

func TestAuthorization_CreateToken_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"code":400,"status":400,"message":"Invalid request"}`)
	})

	_, _, err := client.CreateToken("user", "pass")
	if err == nil {
		t.Fatal("Error expected")

	}
	got := err
	exp := &ErrorResponse{Code: 400, Status: 400, Message: "Invalid request"}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestAuthorization_CreateToken_Success(t *testing.T) {
	setup()
	defer teardown()

	var (
		gotMethod string
	)

	mux.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		fmt.Fprintf(w, `{"token": "secrettoken"}`)
	})

	tok, _, err := client.CreateToken("user", "pass")
	if err != nil {
		t.Fatal(err)
	}

	{
		got := tok
		exp := "secrettoken"
		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}

	{
		got := gotMethod
		exp := "POST"
		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}

}
