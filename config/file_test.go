package config

import "testing"

func TestFileConfigProvider(t *testing.T) {
	var fileConfigProvider = FileConfigProvider{ Path: "./testdata/config-test.yaml" }
	var config = fileConfigProvider.GetConfig()
	if(len(config.Deployments) == 0) {
		t.Errorf("Couldn't read the config file from disk")
	}
}
