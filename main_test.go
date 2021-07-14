package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetEnv(t *testing.T) {
	value := GetEnv("NOT_PRESENT", "default")
	if value != "default" {
		t.Errorf("GetEnv was incorrect, got: %s, want: %s.", value, "default")
	}
}

func TestRedirectHandler(t *testing.T) {
	redirect_uri := "https://onna.com"
	state := base64.StdEncoding.EncodeToString([]byte(redirect_uri))
	uri := fmt.Sprintf("/redirect?state=%s", state)
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	res := httptest.NewRecorder()
	desiredLocation := fmt.Sprintf("%s?state=%s", redirect_uri, state)

	handler(res, req)

	if res.Code != http.StatusFound {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusFound)
	}

	location, _ := url.QueryUnescape(res.Header().Get("Location"))
	if location != desiredLocation {
		t.Errorf("got status %s but wanted %s", location, desiredLocation)
	}
}
