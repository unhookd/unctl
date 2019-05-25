package tests

import (
	"github.com/unhookd/unctl/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unhookd/unctl/client"
	"github.com/unhookd/unctl/helm"
	"github.com/unhookd/unctl/server"
)

func TestCmdZeroTrustServerWithSha(t *testing.T) {
	helm.Setup(false)

	config.CurrentProvider = config.GithubConfigProvider{
		Client: *config.BuildGithubClientFromEnv(),
		Owner: "unhookd",
		Repo: "test-config-store",
		Ref: "master",
		Path: "config-test.yaml",
	}
	config.LoadConfig()

	endpoint := "https://unctl.local:4443/zero-trust"
	values := client.CreateValues("test", "test-deployment", "8b862e527cf90c34a83f0b349bbb686bef33a2cb", "true", "false")

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
	config.CurrentProvider = config.FileConfigProvider{Path: "../config/testdata/config-test.yaml"}
	config.LoadConfig()

	endpoint := "https://unctl.local:4443/zero-trust"
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
