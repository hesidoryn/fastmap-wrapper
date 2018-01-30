package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_fastmap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(fastmap))
	defer ts.Close()

	bbox := "?bbox=27.616,53.853,27.671,53.886"
	ts.URL += bbox
	res1, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	newAPIResponse, err := ioutil.ReadAll(res1.Body)
	res1.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("newApiResponse.xml", newAPIResponse, 0644)

	oldAPIURL := "http://localhost:3000/api/0.6/map" + bbox
	res2, err := http.Get(oldAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	oldAPIResponse, err := ioutil.ReadAll(res2.Body)
	res2.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("oldApiResponse.xml", oldAPIResponse, 0644)

	if bytes.Compare(oldAPIResponse, newAPIResponse) == -1 {
		t.Error("Expected true, got false")
	}
}
