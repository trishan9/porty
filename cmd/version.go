package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the porty version",
	Run: func(cmd *cobra.Command, args []string) {
		showBanner()
		fmt.Println("porty version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
