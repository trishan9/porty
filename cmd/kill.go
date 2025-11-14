package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trishan9/porty/internal"
)

var ports string
var pids string

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill processes by port or PID",
	RunE: func(cmd *cobra.Command, args []string) error {
		showBanner()
		
		if ports == "" && pids == "" {
			return fmt.Errorf("you must specify --port or --pid")
		}

		entries, _ := internal.ListPorts()

		if pids != "" {
			pidList := internal.ParseCSVInts(pids)
			msgs := internal.KillPIDs(pidList)
			for _, m := range msgs {
				fmt.Println(m)
			}
			return nil
		}

		if ports != "" {
			portList := strings.Split(ports, ",")
			msgs := internal.KillByPorts(entries, portList)
			for _, m := range msgs {
				fmt.Println(m)
			}
		}
		return nil
	},
}

func init() {
	killCmd.Flags().StringVar(&ports, "port", "", "Ports to kill (comma-separated)")
	killCmd.Flags().StringVar(&pids, "pid", "", "PIDs to kill (comma-separated)")
	killCmd.Example = `
		porty kill 3000
		porty kill --force 8080
		porty kill --pid 1234
		porty kill 3000 8081 9090
		`
	rootCmd.AddCommand(killCmd)
}
