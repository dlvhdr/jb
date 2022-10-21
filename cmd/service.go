/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"jobber/data"
	"log"
	"os"
	"os/exec"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service started...")
		l := log.New(os.Stderr, "", 0)

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		pubsub := rdb.Subscribe("jobber-new-job")
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage()
			if err != nil {
				panic(err)
			}

			fmt.Printf("New job receieved: %v, %v\n", msg.Channel, msg.Payload)
			var job data.Job
			json.Unmarshal([]byte(msg.Payload), &job)
			if err := job.Parse(); err != nil {
				l.Println(err)
			}

			jobCmd := exec.Command(job.Name, job.Args...)
			stdOutFile, err := os.Create("./out.txt")
			if err != nil {
				l.Println(err)
			}

			stdErrFile, err := os.Create("./err.txt")
			if err != nil {
				l.Println(err)
			}

			jobCmd.Stdout = stdOutFile
			jobCmd.Stderr = stdErrFile
			fmt.Printf("Starting command %+v\n", job)
			err = jobCmd.Start()
			if err != nil {
				l.Println(err)
			} else {
				job.Pid = jobCmd.Process.Pid
				fmt.Printf("Job %v has pid %v\n", job.Id, job.Pid)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
