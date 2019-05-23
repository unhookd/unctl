package providers

import (
	"github.com/unhookd/unctl/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type File struct {
	Path string
	config config.Config
	rawConfig string
}

func (provider File) GetConfig() config.Config {
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

func (provider File) unmarshalConfig() config.Config {
	var config config.Config
	err := yaml.Unmarshal([]byte(provider.rawConfig), &config)

	if err != nil {
		log.Fatalf("Config file is not properly formatted: %v", err)
	}

	return config
}