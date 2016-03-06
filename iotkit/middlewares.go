package iotkit

import (
	"net/http"
	"net/http/httputil"
	// "net/http/httputil"

	"log"
)

type Dumper struct{ Doer }

func (m *Dumper) Do(req *http.Request) (resp *http.Response, err error) {

	dump, err := httputil.DumpRequest(req, true)
	log.Println("Request:", string(dump))

	resp, err = m.Doer.Do(req)
	if err != nil {
		return
	}

	dump, err = httputil.DumpResponse(resp, true)
	log.Println("Response:", string(dump))

	return
}
