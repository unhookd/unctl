package server

import (
	"fmt"
	"github.com/unhookd/unctl/config"
	"net/http"
	"time"

	"github.com/unhookd/unctl/helm"
	"github.com/unhookd/unctl/lib"
)

func ZeroTrustServer(runHelm bool) *http.Server {
	helm.Setup(runHelm)

	http.HandleFunc("/zero-trust", ZeroTrustServerHandler)
	http.HandleFunc("/status", healthCheck)
	listen := "0.0.0.0:8081"
	fmt.Println("Starting unhookD zero-trust... Please standby", listen)

	srv := &http.Server{
		Addr:         listen,
		ReadTimeout:  601 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	srv.ListenAndServe()

	return srv
}

func ZeroTrustServerHandler(w http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(w, request.Body, MaxFileSize)

	var validKeys = []string{"sha", "project", "release", "async", "dryRun"}

	if err := request.ParseForm(); err != nil {
		http.Error(w, err.Error(), 422)
		return
	}

	var allParamKeys []string

	for key := range request.PostForm {
		allParamKeys = append(allParamKeys, key)
	}

	extraParams := lib.Difference(allParamKeys, validKeys)

	if len(extraParams) > 0 {
		w.WriteHeader(201)
		return
	}

	desiredSha := request.FormValue("sha")
	project := request.FormValue("project")
	release := request.FormValue("release")
	async := request.FormValue("async")
	dryRun := request.FormValue("dryRun")

	var effectiveRelease, effectiveNamespace, repo, branch, chart, version string

	var notifications config.NotificationsLookup

	if lookedupProject, ok := config.Current.Deployments[project]; ok {
		if lookedupRelease, ok := lookedupProject[release]; ok {
			effectiveRelease = lookedupRelease.Release
			effectiveNamespace = lookedupRelease.Namespace
			repo = lookedupRelease.Repo
			branch = lookedupRelease.Branch
			chart = lookedupRelease.Chart
			version = lookedupRelease.Version
			notifications = lookedupRelease.Notifications
		}
	}

	if len(effectiveRelease) == 0 || len(effectiveNamespace) == 0 || len(repo) == 0 || len(branch) == 0 || len(chart) == 0 {
		http.Error(w, fmt.Sprintf("Unknown/incomplete project or release: project:%v, release:%v", project, release), 422)
		return
	}

	shouldWait := !lib.VerifyQueryStringParams(async)
	shouldDryRun := lib.VerifyQueryStringParams(dryRun)

	shaToDeploy, validationError := GetShaToDeploy(repo, branch, desiredSha)
	if validationError != nil {
		http.Error(w, fmt.Sprintf("Error validating github args: %v", validationError.Error()), 422)
		return
	}

	values, _ := buildValuesYaml(project, shaToDeploy)
	helmOut, shouldNotify, err := helm.HelmBinUpgradeInstall(effectiveRelease, effectiveNamespace, chart, version, string(values), shouldWait, shouldDryRun)

	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Println(string(values))
	fmt.Println(string(helmOut))
	fmt.Println("++++++++++++++++++++++++++++++++")

	if err != nil {
		http.Error(w, fmt.Sprintf("Error running helm upgrade: %v %v", string(helmOut), err.Error()), 500)
		return
	}

	w.WriteHeader(201)
	w.Write(helmOut)

	fmt.Println("Notifying communication channels")
	if shouldNotify {
		for _, notificationConfigured := range notifications {
			config.NotifyCommunicationChannel(notificationConfigured)
		}
	}
}

func GetShaToDeploy(repo string, branch string, desiredSha string) (sha string, err error) {
	githubClient := config.BuildGithubClientFromEnv()

	headSha, validationError := config.GetHeadSha(repo, branch, githubClient)
	if validationError != nil {
		return "", validationError
	}

	if len(desiredSha) > 0 {
		validationError = config.ValidateShasMatch(headSha, desiredSha)
		if validationError != nil {
			return "", validationError
		}
	}

	validationError = config.ValidateStatusChecks(repo, branch, headSha, githubClient)
	if validationError != nil {
		return "", validationError
	}

	return headSha, nil
}
