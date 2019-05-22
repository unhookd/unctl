package providers

import "testing"
import "github.com/org/unhookd/auth"

func TestGithubConfig(t *testing.T) {
	client := auth.BuildGithubClient()

	var configProvider = Github{
		Client: *client,
		Owner: "unhookd",
		Repo: "test",
		Ref: "master",
		Path: "test.config",
	}
	var config = configProvider.GetConfig()
	if(len(config.Deployments) == 0) {
		t.Errorf("Couldn't read the config file from disk")
	}
}
