package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/unhookd/unctl/config"
	"os"
)

func init() {
	rootCmd.AddCommand(CmdDebugLookup)

	config_provider := os.Getenv("CONFIG_PROVIDER")
	switch config_provider {
	case "github":
		client := config.BuildGithubClientFromEnv()
		config.CurrentProvider = config.GetGithubProviderFromPath(*client, os.Getenv("GITHUB_CONFIG_PATH"))
	case "file":
		config.CurrentProvider = config.FileConfigProvider{Path: os.Getenv("CONFIG_FILE")}
	default:
		config.CurrentProvider = config.FileConfigProvider{Path: "config.yaml"}
	}
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
			for project, releases := range config.Current.Deployments {
				for release, _ := range releases {
					fmt.Println(fmt.Sprintf("/usr/bin/unctl deploy %s %s", project, release))
				}
			}
		} else if len(args) == 1 {
			for release, releaseConfig := range config.Current.Deployments[args[0]] {
				fmt.Println(fmt.Sprintf("/usr/bin/unctl deploy %s %s -- %s@%s", args[0], release, releaseConfig.Repo, releaseConfig.Branch))
			}
		} else if len(args) == 2 {
			releaseConfig := config.Current.Deployments[args[0]][args[1]]
			fmt.Println(fmt.Sprintf("/usr/bin/unctl deploy %s %s\nrepo: %s\nbranch: %s\ncluster: %s\nchart: %s", args[0], args[1], releaseConfig.Repo, releaseConfig.Branch, releaseConfig.Cluster, releaseConfig.Chart))
			fmt.Println(fmt.Sprintf("%v", releaseConfig.Notifications))
		}
	},
}
