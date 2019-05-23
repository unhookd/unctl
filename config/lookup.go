package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Current Config
var CurrentProvider ProviderInterface

type EndpointsTable map[string]string

type ContextsTable map[string][]string

type NotificationsTable map[string]string
type NotificationsLookup []NotificationsTable

type TargetTable struct {
	Release       string
	Repo          string
	Namespace     string
	Cluster       string
	Branch        string
	Chart         string
	Version       string
	Notifications NotificationsLookup
}

type ProjectTable map[string]TargetTable
type DeploymentsTable map[string]ProjectTable

type Config struct {
	Contexts    ContextsTable
	Endpoints   EndpointsTable
	Deployments DeploymentsTable
}

func init() {
	config_provider := os.Getenv("CONFIG_PROVIDER")
	switch config_provider {
	case "github":
		client := BuildGithubClientFromEnv()
		CurrentProvider = GetGithubProviderFromPath(*client, os.Getenv("GITHUB_CONFIG_PATH"))
	case "file":
		CurrentProvider = FileConfigProvider{Path: os.Getenv("CONFIG_FILE")}
	default:
		CurrentProvider = FileConfigProvider{Path: "config.yaml"}
	}
}

func LoadConfig() {
	Current = CurrentProvider.GetConfig()
}

var CmdDebugLookup = &cobra.Command{
	Use:   "config [project] [release]",
	Short: "Shows information about a release",
	Long: `
		Prints information about the config config table
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			for project, releases := range Current.Deployments {
				for release, _ := range releases {
					fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s", project, release))
				}
			}
		} else if len(args) == 1 {
			for release, releaseConfig := range Current.Deployments[args[0]] {
				fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s -- %s@%s", args[0], release, releaseConfig.Repo, releaseConfig.Branch))
			}
		} else if len(args) == 2 {
			releaseConfig := Current.Deployments[args[0]][args[1]]
			fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s\nrepo: %s\nbranch: %s\ncluster: %s\nchart: %s", args[0], args[1], releaseConfig.Repo, releaseConfig.Branch, releaseConfig.Cluster, releaseConfig.Chart))
			fmt.Println(fmt.Sprintf("%v", releaseConfig.Notifications))
		}
	},
}
