package providers

import (
	"github.com/unhookd/unctl/auth"
	"testing"
)

func TestGithubConfig(t *testing.T) {
	//c66ce44344a130c9d513321219a0d120e44afdeb
	// new: 0463bc1ae012a663c9bafff9b59389e5f24f2c97

	client := auth.BuildGithubClient("c66ce44344a130c9d513321219a0d120e44afdeb")

	var configProvider = Github{
		Client: *client,
		Owner: "unhookd",
		Repo: "test-config-store",
		Ref: "master",
		Path: "config-test.yaml",
	}
	var config = configProvider.GetConfig()
	if(len(config.Deployments) == 0) {
		t.Errorf("Couldn't read the config file from disk")
	}
}
