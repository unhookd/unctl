package providers

import (
	"github.com/org/unhookd/auth"
	"testing"
)

func TestGithubConfig(t *testing.T) {
	//c66ce44344a130c9d513321219a0d120e44afdeb

	client := auth.BuildGithubClientFromEnv()

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
