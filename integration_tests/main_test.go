package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unhookd/unctl/client"
	"github.com/unhookd/unctl/helm"
	"github.com/unhookd/unctl/server"
)

func TestCmdZeroTrustServerWithSha(t *testing.T) {
	helm.Setup(false)
	endpoint := "https://local.unhookd.org.net:4443/zero-trust"
	values := client.CreateValues("test", "test-deployment", "adb77bea1a1e80e8da839caa6818b7c56cc8e5b7", "true", "false")

	request, err := http.NewRequest("POST", endpoint, values)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ZeroTrustServerHandler)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	t.Log(responseRecorder.Body.String())
}

func TestCmdZeroTrustServerNoSha(t *testing.T) {
	helm.Setup(false)
	endpoint := "https://local.unhookd.org.net:4443/zero-trust"
	values := client.CreateValues("test", "test-deployment", "", "true", "false")

	request, err := http.NewRequest("POST", endpoint, values)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ZeroTrustServerHandler)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	t.Log(responseRecorder.Body.String())
}
