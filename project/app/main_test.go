package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	expect := Response{Message: "Welcome to the Go REST API!"}

	ts := httptest.NewServer(SetupServer())

	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/", ts.URL))

	if err != nil {
		t.Errorf("Expected no error, got %s\n", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d\n", res.StatusCode)
	}

	var resData Response
	byteData, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(byteData, &resData)

	if err != nil {
		t.Errorf("Failed parsing json, got %s\n", err.Error())
	}

	if resData != expect {
		t.Fatalf("Expected root message, got %v\n", resData)
	}
}
