package providers

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/org/unhookd/config"
	"gopkg.in/yaml.v2"
	"log"
)

type Github struct {
	Client github.Client
	Owner string
	Repo string
	Ref string
	Path string
	config config.Config
	rawConfig string
}

func (provider Github) GetConfig() config.Config {
	ctx := context.Background()
	opts := &github.RepositoryContentGetOptions{ Ref: provider.Ref }
	file_content, _, _, err := provider.Client.Repositories.GetContents(ctx, provider.Owner, provider.Repo, provider.Path, opts)

	if err != nil {
		log.Fatalf("Couldn't find config on Github: %v", err)
	}

	provider.rawConfig, err = file_content.GetContent()

	if err != nil {
		log.Fatalf("Couldn't read config from Github: %v", err)
	}

	return provider.unmarshalConfig()
}

func (provider Github) unmarshalConfig() config.Config {
	var config config.Config
	err := yaml.Unmarshal([]byte(provider.rawConfig), &config)

	if err != nil {
		log.Fatalf("Config file is not properly formatted: %v", err)
	}

	return config
}