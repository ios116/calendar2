package cmd

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/ios116/calendar/internal/config"
	"github.com/ios116/calendar/internal/grpcserver"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

var GrpcClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run grpc client",
	Run: func(cmd *cobra.Command, args []string) {
		client()
	},
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.Token,
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return false
}

func client() {
	option := grpc.WithPerRPCCredentials(&tokenAuth{"Bearer secret"})
	container := BuildContainer()

	container.Invoke(func(conf *config.GrpcConf, logger *zap.Logger) {
		sugar := logger.Sugar()
		address := fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)

		conn, err := grpc.Dial(address, option, grpc.WithInsecure())
		if err != nil {
			sugar.Fatal("Can't connect to GRPC: ", address)
		}

		calendarGRPC := grpcserver.NewCalendarClient(conn)
		ctx := context.Background()

		date, err := ptypes.TimestampProto(time.Now().UTC())
		if err != nil {
			sugar.Fatal(err.Error())
		}
		event := &grpcserver.Event{
			Title:       "My event",
			Date:        date,
			Duration:    ptypes.DurationProto(time.Duration(45 * time.Minute)),
			Author:      "Ivan client",
			Description: "Some descriptions",
			Notify:      ptypes.DurationProto(time.Duration(60 * time.Hour)),
		}
		// add event
		eventResp, err := calendarGRPC.CreateEvent(ctx, event)
		if err != nil {
			sugar.Fatal(err.Error())
		}
		sugar.Info(eventResp)

		// edit event
		event.Id = eventResp.Event.Id
		event.Title = "Edited event 45"
		respFromEdit, err := calendarGRPC.UpdateEvent(ctx, event)
		if err != nil {
			logger.Fatal(err.Error())
		}
		sugar.Info(respFromEdit.Detail)

		// get by period

		oldDate := time.Now().UTC().Add(-time.Duration(24 * time.Hour))
		oldDateProto, _ := ptypes.TimestampProto(oldDate)
		per := grpcserver.PeriodRequest{
			Period: grpcserver.Periods_WEEK,
			Date:   oldDateProto,
		}

		evID := &grpcserver.EventIDRequest{
			Id: 55,
		}
		respDel, err := calendarGRPC.DeleteEvent(ctx, evID)
		if err != nil {
			logger.Fatal(err.Error())
		}
		sugar.Info(respDel)

		respFromPeriod, err := calendarGRPC.GetEvents(ctx, &per)
		if err != nil {
			sugar.Fatal(err.Error())
		}
		sugar.Info(respFromPeriod.Events)
	})

}
