package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

var GlobalLookups Config
var EncodedConfigLookup string

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

// Probably don't need this anymore
//func init() {
//	if len(EncodedConfigLookup) == 0 {
//		log.Fatalf("missing config, ensure bin is built with proper ldflag")
//	}
//
//	decoded, err := base64.StdEncoding.DecodeString(EncodedConfigLookup)
//	if err != nil {
//		log.Fatalf("decode error: %v", err)
//	}
//
//	err = yaml.Unmarshal([]byte(decoded), &GlobalLookups)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//}

var CmdDebugLookup = &cobra.Command{
	Use:   "config [project] [release]",
	Short: "Shows information about a release",
	Long: `
		Prints information about the config config table
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
