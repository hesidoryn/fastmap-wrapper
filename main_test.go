package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_fastmapBadRequests(t *testing.T) {
	t.Run("empty bbox", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=0,0,0", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=0,0,0")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=0.1,0.1,-0.1,-0.1", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=0.1,0.1,-0.1,-0.1")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=179.9,0,180.1,0.1", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=179.9,0,180.1,0.1")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=-180.1,0,-179.9,0.1", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=-180.1,0,-179.9,0.1")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=0,89.9,0.1,90.1", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=0,89.9,0.1,90.1")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("bbox=0,-90.1,0.1,-89.9", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(fastmap))
		defer ts.Close()
		res, err := http.Get(ts.URL + "?bbox=0,-90.1,0.1,-89.9")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, res.StatusCode)
		}
	})
}

func Test_fastmap(t *testing.T) {
	t.Run("bbox - table test", func(t *testing.T) {
		bboxes := []string{
			"27.616,53.853,27.630,53.870",
			"27.616,53.853,27.617,53.854",
			"27.61649608612061,53.85379229563698,27.671985626220707,53.886459293813054",
		}
		for i := range bboxes {
			path := "xml_files/" + bboxes[i]
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.MkdirAll(path, os.ModePerm)
			}

			// ts := httptest.NewServer(http.HandlerFunc(fastmap))
			bbox := "?bbox=" + bboxes[i]
			// ts.URL += bbox
			fastmapURL := "http://localhost:3001/api/0.6/map" + bbox
			startReq := time.Now()
			// res1, err := http.Get(ts.URL)
			res1, err := http.Get(fastmapURL)
			if err != nil {
				log.Fatal(err)
			}
			t.Log(bbox, "fastmap response time: ", time.Since(startReq))
			fastmapResponse, err := ioutil.ReadAll(res1.Body)
			res1.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			ioutil.WriteFile(path+"/fastmap.xml", fastmapResponse, 0644)
			fm := &osm{}
			xml.Unmarshal(fastmapResponse, fm)
			fm.sort()

			portURL := "http://localhost:3000/api/0.6/map" + bbox
			startReq = time.Now()
			res2, err := http.Get(portURL)
			if err != nil {
				log.Fatal(err)
			}
			t.Log(bbox, "port response time: ", time.Since(startReq))
			portResponse, err := ioutil.ReadAll(res2.Body)
			res2.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			ioutil.WriteFile(path+"/port.xml", portResponse, 0644)
			port := &osm{}
			xml.Unmarshal(portResponse, port)
			port.sort()

			if !reflect.DeepEqual(fm, port) {
				t.Errorf("bbox=%v, expected: true, got: false", bboxes[i])
			}

			cgimapURL := "http://localhost:31337/api/0.6/map" + bbox
			startReq = time.Now()
			res3, err := http.Get(cgimapURL)
			if err != nil {
				log.Fatal(err)
			}
			t.Log(bbox, "cgimap response time: ", time.Since(startReq))
			cgimapResponse, err := ioutil.ReadAll(res3.Body)
			res3.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			ioutil.WriteFile(path+"/cgimap.xml", cgimapResponse, 0644)
			cgimap := &osm{}
			xml.Unmarshal(cgimapResponse, cgimap)
			cgimap.sort()

			if !reflect.DeepEqual(fm, cgimap) {
				t.Errorf("bbox=%v, expected: true, got: false", bboxes[i])
			}

			// ts.Close()
		}
	})
}
