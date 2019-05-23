package config

import (
	"testing"
)

func TestGetHeadSha(t *testing.T) {
	CurrentProvider = FileConfigProvider{ Path: "./testdata/config-test.yaml" }
	LoadConfig()

	repo, branch := "", ""
	if lookedupProject, ok := Current.Deployments["test"]; ok {
		if lookedupRelease, ok := lookedupProject["test-deployment"]; ok {
			repo = lookedupRelease.Repo
			branch = lookedupRelease.Branch
		}
	}
	client := BuildGithubClientFromEnv()
	sha, _ := GetHeadSha(repo, branch, client)

	if len(sha) == 0 {
		t.Errorf("Failed to return the head sha of branch %s in repo %s", branch, repo)
	}

}
