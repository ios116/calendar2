package grpc_test

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"integration_tests/internal/config"
	"integration_tests/internal/grpcserver"
	"log"
	"testing"
	"time"
)

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
func TestGRPC(t *testing.T) {

	option := grpc.WithPerRPCCredentials(&tokenAuth{"Bearer secret"})
	conf := config.NewGrpcConf()
	address := fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)
	conn, err := grpc.Dial(address, option, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Can't connect to GRPC: ", address)
	}

	calendarGRPC := grpcserver.NewCalendarClient(conn)
	ctx := context.Background()

	date, err := ptypes.TimestampProto(time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	event := &grpcserver.Event{
		Id:          0,
		Title:       "My event",
		Date:        date,
		Duration:    ptypes.DurationProto(time.Duration(45 * time.Minute)),
		Author:      "Ivan client",
		Description: "Some descriptions",
		Notify:      ptypes.DurationProto(time.Duration(60 * time.Hour)),
	}

	t.Run("add event", func(t *testing.T) {
		eventResp, err := calendarGRPC.CreateEvent(ctx, event)
		if err != nil {
			t.Log(err)
		}
		event = eventResp.Event
	})

	t.Run("edit event", func(t *testing.T) {
		event.Title = "Edited event"
		_, err := calendarGRPC.UpdateEvent(ctx, event)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Get by id", func(t *testing.T) {
		t.Log("event id=",event.Id)
		req := grpcserver.EventIDRequest{
			Id:                   event.Id,
		}
	eventById, err:= calendarGRPC.GetEvent(ctx, &req)
	if err !=nil {
	   t.Fatal(err)
	}
	if eventById.Event.Id != event.Id {
		t.Fatal("id is not equal")
	}
	})

	t.Run("Delete by id", func(t *testing.T) {
		req := grpcserver.EventIDRequest{
			Id:                   event.Id,
		}
		_, err:= calendarGRPC.DeleteEvent(ctx, &req)
		if err !=nil {
			t.Fatal(err)
		}
	})

	t.Run("get by period", func(t *testing.T) {
		oldDate := time.Now().UTC().Add(-time.Duration(24 * time.Hour))
		oldDateProto, _ := ptypes.TimestampProto(oldDate)
		per := grpcserver.PeriodRequest{
			Period: grpcserver.Periods_WEEK,
			Date:   oldDateProto,
		}
		respFromPeriod, err := calendarGRPC.GetEvents(ctx, &per)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(len(respFromPeriod.Events))
	})
}
