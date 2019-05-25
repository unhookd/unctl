package config

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifySlack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	channel := "#ops-notes"
	message := "foobar"
	payload := fmt.Sprintf(`{"channel": "%v", "username": "unctl", "text": "%v", "icon_emoji": ":ops:"}`, channel, message)
	sn, response, _ := notifySlack(ts.URL, channel, message)

	if payload != response {
		t.Errorf("The payload did not match the response")
	}

	if !sn {
		t.Errorf("Failed to post to slack.")
	}
}

func TestNotifySlack_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	channel := "#integrations-platform"
	message := "foobar"
	payload := fmt.Sprintf(`{"channel": "%v", "username": "unctl", "text": "%v", "icon_emoji": ":ops:"}`, channel, message)
	sn, response, _ := notifySlack(ts.URL, channel, message)

	if payload != response {
		t.Errorf("The payload did not match the response")
	}

	if sn {
		t.Errorf("Sucessfully posted to slack.")
	}
}
