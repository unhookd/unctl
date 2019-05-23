package config

import (
	"github.com/unhookd/unctl/auth"
	"testing"
)

func TestGithubConfigProvider(t *testing.T) {
	client := auth.BuildGithubClientFromEnv()

	var configProvider = GithubConfigProvider{
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
