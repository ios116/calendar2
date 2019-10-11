package cmd

import (
	"context"
	"github.com/ios116/calendar/internal/web"
	"github.com/spf13/cobra"
	"log"
)

var HttpServerCmd = &cobra.Command{
	Use:   "web",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	container := BuildContainer()
	err := container.Invoke(func(webServer *web.HttpServer) {
		webServer.Run()
	})
	if err != nil {
		log.Println(err)
	}

}
