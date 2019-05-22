package auth

import (
	"testing"
	"github.com/org/unhookd/config"
)

func TestGetHeadSha(t *testing.T) {
	repo, branch := "", ""
	if lookedupProject, ok := lookup.GlobalLookups.Deployments["test"]; ok {
		if lookedupRelease, ok := lookedupProject["test-deployment"]; ok {
			repo = lookedupRelease.Repo
			branch = lookedupRelease.Branch
		}
	}
	client := BuildGithubClient()
	sha, _ := GetHeadSha(repo, branch, client)

	if len(sha) == 0 {
		t.Errorf("Failed to return the head sha of branch %s in repo %s", branch, repo)
	}

}
