package main

import (
	"github.com/unhookd/unctl/client"
	"github.com/unhookd/unctl/server"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{Use: "unhookd"}
	rootCmd.PersistentFlags().String("debug", "false", "whether or not to show debug logic")
	rootCmd.AddCommand(client.CmdDeploy)
	rootCmd.AddCommand(server.CmdZeroTrustServer)
	rootCmd.AddCommand(config.CmdDebugLookup)
	rootCmd.Execute()

}
