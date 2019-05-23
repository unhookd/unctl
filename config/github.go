package config

import (
	"context"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type GithubConfigProvider struct {
	Client    github.Client
	Owner     string
	Repo      string
	Ref       string
	Path      string
	config    Config
	rawConfig string
}

func GetGithubProviderFromPath(client github.Client, fullPath string) GithubConfigProvider {
	parts := strings.Split(fullPath, "/")
	owner, parts := parts[0], parts[1:]
	repo, parts := parts[0], parts[1:]
	ref, parts := parts[0], parts[1:]
	path := strings.Join(parts, "/")

	return GithubConfigProvider{
		Client: client,
		Owner: owner,
		Repo: repo,
		Ref: ref,
		Path: path,
	}
}

func (provider GithubConfigProvider) GetConfig() Config {
	ctx := context.Background()
	opts := &github.RepositoryContentGetOptions{ Ref: provider.Ref }
	file_content, _, _, err := provider.Client.Repositories.GetContents(ctx, provider.Owner, provider.Repo, provider.Path, opts)

	if err != nil {
		log.Fatalf("Couldn't find config on GithubConfigProvider: %v", err)
	}

	provider.rawConfig, err = file_content.GetContent()

	if err != nil {
		log.Fatalf("Couldn't read config from GithubConfigProvider: %v", err)
	}

	return provider.unmarshalConfig()
}

func (provider GithubConfigProvider) unmarshalConfig() Config {
	var config Config
	err := yaml.Unmarshal([]byte(provider.rawConfig), &config)

	if err != nil {
		log.Fatalf("Config file is not properly formatted: %v", err)
	}

	return config
}