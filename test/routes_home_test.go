package test

import (
	"metasource/metasource/option"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome_UnInit(t *testing.T) {
	defer Path_UnInit("/var/tmp")()

	router := option.Navigate()
	rqster := httptest.NewRequest("GET", "/", nil)
	record := httptest.NewRecorder()
	router.ServeHTTP(record, rqster)

	subslist := []string{
		"MetaSource is a performant source for RPM repositories metadata which ",
		"has an access to the metadata of the different Fedora Linux package ",
		"repositories and will serve you the most recent information available.",
		"It will parse through the \"updates-testing\" repository before moving",
		"onto the likes of \"updates\" and \"releases\" repository if no ",
		"information is found in the previous repository.",
		"Utilize the fast lookup interface to acquaint yourself with the API ",
		"endpoints and expected outputs. Press ENTER after typing the name to ",
		"execute a lookup in a new window. If you query for a non-existent ",
		"branch - it will return an HTTP 400 error. If you query for a ",
		"non-existent package - it will return an HTTP 404 error. Please report",
		"persistent HTTP 500 errors to the Fedora Infrastructure team.",
		"You can retrieve the information about available branches by querying ",
		"the following request.",
		"Please ensure that the database directory has been configured properly",
		"and the scheduled downloads have been functioning correctly.",
	}
	for _, item := range subslist {
		if !strings.Contains(record.Body.String(), item) {
			t.Errorf("Absent '%s'", item)
		}
	}

	if !strings.Contains(record.Header().Get("content-type"), "text/html") {
		t.Errorf("Received %s, Expected %s", record.Header().Get("content-type"), "text/html")
	}

	if record.Code != http.StatusOK {
		t.Errorf("Received %d, Expected %d", record.Code, http.StatusOK)
	}
}
