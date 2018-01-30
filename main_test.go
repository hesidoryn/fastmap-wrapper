package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_fastmap(t *testing.T) {
	t.Run("Bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=13123234142141241234")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox - table test", func(t *testing.T) {
		bboxes := []string{"27.616,53.853,27.630,53.870", "27.616,53.853,27.617,53.854"}
		for i := range bboxes {
			ts := httptest.NewServer(http.HandlerFunc(fastmap))
			bbox := "?bbox=" + bboxes[i]
			ts.URL += bbox
			res1, err := http.Get(ts.URL)
			if err != nil {
				log.Fatal(err)
			}
			fastmapResponse, err := ioutil.ReadAll(res1.Body)
			res1.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			fm := &osm{}
			xml.Unmarshal(fastmapResponse, fm)
			fm.sort()

			portURL := "http://localhost:3000/api/0.6/map" + bbox
			res2, err := http.Get(portURL)
			if err != nil {
				log.Fatal(err)
			}
			portResponse, err := ioutil.ReadAll(res2.Body)
			res2.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			port := &osm{}
			xml.Unmarshal(portResponse, port)
			port.sort()

			if !reflect.DeepEqual(fm, port) {
				ioutil.WriteFile("fastmap-"+bboxes[i]+".xml", fastmapResponse, 0644)
				ioutil.WriteFile("port-"+bboxes[i]+".xml", fastmapResponse, 0644)
				t.Errorf("bbox=%v, expected: true, got: false", bboxes[i])
			}

			ts.Close()
		}
	})
}
