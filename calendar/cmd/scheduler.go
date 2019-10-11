package cmd

import (
	"github.com/ios116/calendar/internal/scheduler"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var RQConsumerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Run RQ scanner",
	Run: func(cmd *cobra.Command, args []string) {
		time.Sleep(10*time.Second)
		container := BuildContainer()
		err := container.Invoke(func(scanner *scheduler.Scanner) {
			scanner.Produce()
		})
		log.Println(err)
	},
}
