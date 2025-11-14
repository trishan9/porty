package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trishan9/porty/internal"
	"github.com/trishan9/porty/tui"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display active ports in an interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		showBanner()
		entries, err := internal.ListPorts()
		if err != nil {
			return err
		}

		if jsonOutput {
			b, err := json.MarshalIndent(entries, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		}

		if err := tui.Run(entries); err != nil {
			fmt.Fprintln(os.Stderr, "TUI error:", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
