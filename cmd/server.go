package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unhookd/unctl/server"
)

var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Runs the unctl server",
	Long: `
	  Runs the unctl server, exposing an endpoint where deploys defined in the unctl server config can be triggered.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		server.ZeroTrustServer(true)
	},
}

func init() {
	rootCmd.AddCommand(CmdServer)
}