package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type FileConfigProvider struct {
	Path      string
	config    Config
	rawConfig string
}

func (provider FileConfigProvider) GetConfig() Config {
	dat, err := ioutil.ReadFile(provider.Path)
	provider.rawConfig = string(dat)

	if err != nil {
		log.Fatalf("Couldn't read config from filesystem: %v", err)
	}

	if len(provider.rawConfig) == 0 {
		log.Fatalf("Config is empty!")
	}

	return provider.unmarshalConfig()
}

func (provider FileConfigProvider) unmarshalConfig() Config {
	var config Config
	err := yaml.Unmarshal([]byte(provider.rawConfig), &config)

	if err != nil {
		log.Fatalf("Config file is not properly formatted: %v", err)
	}

	return config
}