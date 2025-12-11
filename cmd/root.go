package cmd

import (
    "fmt"
    "os"

    "github.com/charmbracelet/lipgloss"
    "github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "porty",
	Short: "porty - Aesthetic port manager for Linux",
	Long: `porty is a beautiful, modern port management CLI. 
It shows color-coded ports, lets you filter, kill ports, and export port details to JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
			showBanner()
			cmd.Help()
	},
}

func showBanner() {
	banner := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#7dcfff")).
    Render(PortyBanner)
	fmt.Println(banner)
}

var PortyBanner = `
██████╗  ██████╗ ██████╗ ████████╗██╗   ██╗
██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝╚██╗ ██╔╝
██████╔╝██║   ██║██████╔╝   ██║    ╚████╔╝ 
██╔     ██║   ██║██╔══██╗   ██║     ╚██╔╝  
██╔      ██████╔╝██║  ██║   ██║      ██║   
╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝      ╚═╝    
     A modern, and minimal port manager     
`

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}
