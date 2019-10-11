package cmd

import (
	"context"
	"github.com/ios116/calendar/internal/calendar"
	"github.com/ios116/calendar/internal/config"
	"github.com/ios116/calendar/internal/domain"
	"github.com/ios116/calendar/internal/grpcserver"
	"github.com/ios116/calendar/internal/scheduler"
	"github.com/ios116/calendar/internal/sender"
	"github.com/ios116/calendar/internal/store"
	"github.com/ios116/calendar/internal/web"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "CleanCalendar is a calendar microservice",
}

func init() {
	RootCmd.AddCommand(GrpcServerCmd, GrpcClientCmd, RQConsumerCmd, RQProducerCmd, HttpServerCmd)
}

func DbRepo(s *store.DbRepo) domain.EventRepository {
	return domain.EventRepository(s)
}

func UseCaseCalendar(c *calendar.Calendar) calendar.UseCaseCalendar {
	return calendar.UseCaseCalendar(c)
}

func Mail(s *sender.MailService) sender.Mailer {
	return sender.Mailer(s)
}

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(config.NewAppConf)
	container.Provide(config.NewHttpConf)
	container.Provide(config.NewDateBaseConf)
	container.Provide(config.NewGrpcConf)

	container.Provide(config.CreateLogger)
	container.Provide(config.DBConnection)

	container.Provide(context.Background)
	container.Provide(store.NewDbRepo)
	container.Provide(DbRepo)
	container.Provide(calendar.NewCalendar)
	container.Provide(grpcserver.NewGRPCServer)

	container.Provide(calendar.NewCalendar)
	container.Provide(UseCaseCalendar)

	container.Provide(config.NewRabbitConf)
	container.Provide(config.RQConnection)

	container.Provide(sender.NewMailService)
	container.Provide(Mail)
	container.Provide(sender.NewSender)
	container.Provide(scheduler.NewScanner)
	container.Provide(web.NewHttpServer)

	return container
}
