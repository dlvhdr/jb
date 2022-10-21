/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"jobber/data"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		id := fmt.Sprintf("job-%v", time.Now().Unix())
		job := strings.Join(args, " ")

		err := rdb.Publish("jobber-new-job", data.Job{
			Id:      id,
			Cmdline: job,
		}).Err()
		if err != nil {
			panic(err)
		}

		s := strings.Builder{}
		s.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("2")).
				Bold(true).
				Render("Job added: "),
		)
		s.WriteString(
			lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Padding(0, 1).
				Render(job),
		)
		s.WriteString("\n")
		fmt.Printf(s.String())
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
