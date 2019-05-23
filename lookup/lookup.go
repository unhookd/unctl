package lookup

import (
	"fmt"
	"github.com/spf13/cobra"
)

var GlobalLookups Config

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

var CmdDebugLookup = &cobra.Command{
	Use:   "lookup [project] [release]",
	Short: "Shows information about a release",
	Long: `
		Prints information about the lookup config table
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			for project, releases := range GlobalLookups.Deployments {
				for release, _ := range releases {
					fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s", project, release))
				}
			}
		} else if len(args) == 1 {
			for release, releaseConfig := range GlobalLookups.Deployments[args[0]] {
				fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s -- %s@%s", args[0], release, releaseConfig.Repo, releaseConfig.Branch))
			}
		} else if len(args) == 2 {
			releaseConfig := GlobalLookups.Deployments[args[0]][args[1]]
			fmt.Println(fmt.Sprintf("/usr/bin/unhookd deploy %s %s\nrepo: %s\nbranch: %s\ncluster: %s\nchart: %s", args[0], args[1], releaseConfig.Repo, releaseConfig.Branch, releaseConfig.Cluster, releaseConfig.Chart))
			fmt.Println(fmt.Sprintf("%v", releaseConfig.Notifications))
		}
	},
}
