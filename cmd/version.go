package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of unctl",
	Long:  `Print the version of unctl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unctl v0.9 -- HEAD")
	},
}
