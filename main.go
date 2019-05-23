package main

import (
	"github.com/spf13/cobra"
	"github.com/unhookd/unctl/client"
	"github.com/unhookd/unctl/config"
	"github.com/unhookd/unctl/server"
)

func main() {
	var rootCmd = &cobra.Command{Use: "unhookd"}
	rootCmd.PersistentFlags().String("debug", "false", "whether or not to show debug logic")
	rootCmd.AddCommand(client.CmdDeploy)
	rootCmd.AddCommand(server.CmdZeroTrustServer)
	rootCmd.AddCommand(config.CmdDebugLookup)
	rootCmd.Execute()
}
