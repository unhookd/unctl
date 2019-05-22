package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/org/unhookd/config"
	"github.com/org/unhookd/lib"
	"golang.org/x/oauth2"
	"os"
	"regexp"
	"strings"
)



func BuildGithubClient() (client *github.Client) {
	ctx := context.Background()
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	if githubAccessToken == "" {
		panic(errors.New("No GITHUB_ACCESS_TOKEN env"))
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)

	return client
}

func ValidateShasMatch(headSha string, sha string) (err error) {
	isSha, _ := regexp.MatchString("^[0-9a-f]{40}$", sha)

	if !isSha {
		return errors.New(fmt.Sprintf("Unable to parse sha"))
	}

	if headSha != sha {
		return errors.New(fmt.Sprintf("sha: '%s' is not the head of the branch", sha))
	}

	return nil
}

func GetHeadSha(orgRepo, branch string, client *github.Client) (sha string, err error){
	var org, repo = splitOrgRepo(orgRepo)
	headRef, _, err := client.Git.GetRef(context.Background(), org, repo, "refs/heads/"+branch)
	if headRef == nil {
		return "", err
	}
	return *headRef.Object.SHA, nil
}

func splitOrgRepo(orgRepo string) (string, string) {
	var parts = strings.Split(orgRepo, "/")
	return parts[0], parts[1]
}

func ValidateStatusChecks(orgRepo, branch, headSha string, client *github.Client) (err error) {
	var org, repo = splitOrgRepo(orgRepo)
	statuses, _, err := client.Repositories.ListStatuses(context.Background(), org, repo, headSha, &github.ListOptions{Page: 1})
	if err != nil {
		return err
	}

	var required_contexts []string
	var ok bool
	if required_contexts, ok = config.GlobalLookups.Contexts[orgRepo]; !ok {
		return errors.New(fmt.Sprintf("Unable to locate required_contexts: %s", orgRepo))
	}

	var found []string

	for _, r := range required_contexts {
		for _, s := range statuses {
			if *s.State == "success" && *s.Context == r {
				fmt.Println("name:", *s.Context, *s.Description, *s.CreatedAt, "state:", *s.State)
				found = append(found, *s.Context)
			}
		}
	}

	missingStatusChecks := lib.Difference(required_contexts, found)

	if len(missingStatusChecks) > 0 {
		return errors.New(fmt.Sprintf("Status Checks Missing: %s", missingStatusChecks))
	}

	return nil
}
