package client

import (
	"github.com/unhookd/unctl/config"
	"testing"
)

func TestGetGithubArgs(t *testing.T) {
	config.CurrentProvider = config.FileConfigProvider{ Path: "../config/testdata/config-test.yaml" }
	config.LoadConfig()

	args := []string{"test", "test-deployment", "adb77bea1a1e80e8da839caa6818b7c56cc8e5b7"}
	project, release, sha := getGithubArgs(args)

	if project != args[0] && release != args[1] && sha != args[2] {
		t.Errorf("Expected %v, %v, %v. Actual %v, %v, %v", args[0], args[1], args[2], project, release, sha)
	}
}

func TestZeroTrustLookup(t *testing.T) {
	config.CurrentProvider = config.FileConfigProvider{ Path: "../config/testdata/config-test.yaml" }
	config.LoadConfig()

	args := []string{"test", "test-deployment"}
	project, release, _ := getGithubArgs(args)
	cluster, endpoint := zeroTrustLookup(project, release)
	expectedEndpoint := "http://localhost:8081/zero-trust"
	expectedCluster := "local"

	if cluster != expectedCluster && endpoint != expectedEndpoint {
		t.Errorf("Expected %v, %v. Actual %v, %v", cluster, endpoint, expectedCluster, expectedEndpoint)
	}
}