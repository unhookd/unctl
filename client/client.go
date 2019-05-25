package client

import (
	"fmt"
	"github.com/unhookd/unctl/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ClientDeploy(args []string, noWait, shouldDryRun bool) (string, int) {
	project, release, sha := getGithubArgs(args)

	cluster, endpoint := zeroTrustLookup(project, release)

	if len(project) == 0 || len(release) == 0 || len(cluster) == 0 || len(endpoint) == 0 {
		log.Fatalf("missing argument(s): project:%v, release:%v, cluster:%v, endpoint:%v", project, release, cluster, endpoint)
	}

	fmt.Println(
		fmt.Sprintf("found project:%v, release:%v, sha:%v, cluster:%v, endpoint:%v, nowait:%v, dry-run:%v",
			project, release, sha, cluster, endpoint, noWait, shouldDryRun))

	async := ""
	if noWait {
		async = "true"
	}

	dryRun := ""
	if shouldDryRun {
		dryRun = "true"
	}

	encodedValues := CreateValues(project, release, sha, async, dryRun)

	return makeRequest(endpoint, encodedValues)
}

func getGithubArgs(args []string) (string, string, string) {
	var project = args[0]
	var release = args[1]
	sha := ""
	if len(args) == 3 {
		sha = args[2]
	}
	return project, release, sha
}

func zeroTrustLookup(project string, release string) (string, string) {
	var cluster, endpoint string

	if lookedupProject, ok := config.Current.Deployments[project]; ok {
		if lookedupRelease, ok := lookedupProject[release]; ok {
			cluster = lookedupRelease.Cluster
			endpoint = config.Current.Endpoints[cluster]
		}
	}

	return cluster, endpoint
}

func CreateValues(project string, release string, sha string, async string, dryRun string) *strings.Reader {
	values := url.Values{}
	if len(sha) > 0 {
		values.Add("sha", sha)
	}
	values.Add("project", project)
	values.Add("release", release)
	values.Add("async", async)
	values.Add("dryRun", dryRun)

	return strings.NewReader(values.Encode())
}

func makeRequest(endpoint string, encodedValues *strings.Reader) (string, int) {
	req, err := http.NewRequest("POST", endpoint, encodedValues)
	if err != nil {
		log.Fatalf("Unable to create a valid request: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Unable to send POST request to unctl server: %v", err.Error())
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	bodyString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v", err.Error())
	}

	fmt.Println("response Body:", string(bodyString))

	if resp.StatusCode != 201 {
		log.Fatalf("Unable to deploy")
	}

	return string(bodyString), resp.StatusCode
}
