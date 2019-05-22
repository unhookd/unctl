package client

import (
	"github.com/spf13/cobra"
)

var version string
var noWait bool
var shouldDryRun bool

var CmdDeploy = &cobra.Command{
	Use:   "deploy [project] [release] [sha]",
	Short: "Deploy an authorized project to a release in Kubernetes",
	Long: `
	  Given a [project] [release] and optional [sha], an HTTP request is made to the designated
	  zero-trust-server endpoint, and a request to deploy is processed. If no sha is provided, 
	  the head of the branch specified in the zero-trust.yaml will be deployed.
	`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ClientDeploy(args, noWait, shouldDryRun)
	},
}

func init() {
	CmdDeploy.Flags().StringVarP(&version, "version", "v", "", "what version to deploy")
	CmdDeploy.Flags().BoolVar(&noWait, "no-wait", false, "whether or not the command should wait for the deploy to finish")
	CmdDeploy.Flags().BoolVar(&shouldDryRun, "dry-run", false, "noop deploy (dry-run)")
}
