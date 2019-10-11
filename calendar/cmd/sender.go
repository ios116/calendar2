package cmd

import (
	"github.com/ios116/calendar/internal/sender"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var RQProducerCmd = &cobra.Command{
	Use:   "sender",
	Short: "Run RQ sender",
	Run: func(cmd *cobra.Command, args []string) {
		time.Sleep(10*time.Second)
		container := BuildContainer()
		err := container.Invoke(func(sender *sender.Sender) {
			sender.Consume()
		})
		log.Println(err)
	},
}
