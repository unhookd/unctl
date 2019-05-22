package providers

import "testing"

func TestFileConfig(t *testing.T) {
	var fileConfigProvider = File{ Path: "./testdata/config-test.yaml" }
	var config = fileConfigProvider.GetConfig()
	if(len(config.Deployments) == 0) {
		t.Errorf("Couldn't read the config file from disk")
	}
}
