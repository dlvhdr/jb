/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"jobber/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func createModel() ui.Model {
	return ui.NewModel()
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := createModel()

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
