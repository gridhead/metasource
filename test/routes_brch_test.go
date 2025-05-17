package test

import (
	"metasource/metasource/option"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBranches_UnInit(t *testing.T) {
	defer Path_UnInit("/var/tmp")()

	router := option.Navigate()
	rqster := httptest.NewRequest("GET", "/branches", nil)
	record := httptest.NewRecorder()
	router.ServeHTTP(record, rqster)

	if !strings.Contains(record.Header().Get("content-type"), "application/json") {
		t.Errorf("Received %s, Expected %s", record.Header().Get("content-type"), "application/json")
	}

	if record.Body.String() != "[]\n" {
		t.Errorf("Received %s, Expected %s", record.Body.String(), "[]")
	}

	if record.Code != http.StatusOK {
		t.Errorf("Received %d, Expected %d", record.Code, http.StatusOK)
	}
}

func TestBranches_Init(t *testing.T) {
	defer Path_Init()()

	router := option.Navigate()
	rqster := httptest.NewRequest("GET", "/branches", nil)
	record := httptest.NewRecorder()
	router.ServeHTTP(record, rqster)

	if !strings.Contains(record.Header().Get("content-type"), "application/json") {
		t.Errorf("Received %s, Expected %s", record.Header().Get("content-type"), "application/json")
	}

	if record.Body.String() != "[\"rawhide\"]\n" {
		t.Errorf("Received %s, Expected %s", record.Body.String(), "[\"rawhide\"]")
	}

	if record.Code != http.StatusOK {
		t.Errorf("Received %d, Expected %d", record.Code, http.StatusOK)
	}
}
