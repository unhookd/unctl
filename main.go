package main

import (
	"github.com/org/unhookd/client"
	"github.com/org/unhookd/server"
	"github.com/org/unhookd/lookup"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{Use: "unhookd"}
	rootCmd.PersistentFlags().String("debug", "false", "whether or not to show debug logic")
	rootCmd.AddCommand(client.CmdDeploy)
	rootCmd.AddCommand(server.CmdInstastageServer)
	rootCmd.AddCommand(server.CmdZeroTrustServer)
	rootCmd.AddCommand(lookup.CmdDebugLookup)
	rootCmd.Execute()

}
