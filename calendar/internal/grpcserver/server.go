package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/ios116/calendar/internal/calendar"
	"github.com/ios116/calendar/internal/config"
	"github.com/ios116/calendar/internal/domain"
	"github.com/ios116/calendar/internal/exceptions"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

// GRPCServer - grpc server
type GRPCServer struct {
	logger   *zap.Logger
	calendar calendar.UseCaseCalendar
	conf     *config.GrpcConf
}

// NewGRPCServer - constructor for GRPCServer
func NewGRPCServer(logger *zap.Logger, calendar calendar.UseCaseCalendar, conf *config.GrpcConf) *GRPCServer {
	return &GRPCServer{logger: logger, calendar: calendar, conf: conf}
}

// Start - init RPC server
func (g *GRPCServer) Start() {
	address := fmt.Sprintf("%s:%d", g.conf.GrpcHost, g.conf.GrpcPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		g.logger.Fatal("Cannot start RPC server", zap.String("err", err.Error()))
	}

	// server := grpc.NewServer(grpc.UnaryInterceptor(newInterceptor(g.logger, g.conf.GrpcToken)))
	server := grpc.NewServer()
	RegisterCalendarServer(server, g)
	g.logger.Info("Starting RPC server", zap.String("address", address))
	err = server.Serve(lis)
	if err != nil {
		g.logger.Fatal("Cannot start listen port", zap.String("err", err.Error()))
	}
}

// CreateEvent create a event
func (g *GRPCServer) CreateEvent(ctx context.Context, in *Event) (*EventResponse, error) {
	ev, err := toEvent(in)
	if err != nil {
		return nil, err
	}
	event, err := g.calendar.Add(ev)
	switch err {
	case nil:
		rpcEvent, _ := fromEven(event)
		return &EventResponse{
			Status: true,
			Event:  rpcEvent,
			Detail: "Record has been added",
		}, nil
	default:
		g.logger.Error(err.Error())
		var dErr exceptions.DomainError
		if errors.As(err, &dErr) {
			return &EventResponse{
				Status: false,
				Event:  nil,
				Detail: err.Error(),
			}, nil
		} else {
			return nil, err
		}
	}
}

// UpdateEvent update a event
func (g *GRPCServer) UpdateEvent(ctx context.Context, in *Event) (*StatusResponse, error) {

	event, err := toEvent(in)
	st, err := g.calendar.Edit(event)
	switch err {
	case nil:
		return &StatusResponse{
			Status: st,
			Detail: "Record has been edited",
		}, nil
	default:
		g.logger.Error(err.Error())
		var dErr exceptions.DomainError
		if errors.As(err, &dErr) {
			return &StatusResponse{
				Status: false,
				Detail: err.Error(),
			}, nil
		} else {
			return nil, err
		}
	}
}

// DeleteEvent a event
func (g *GRPCServer) DeleteEvent(ctx context.Context, in *EventIDRequest) (*StatusResponse, error) {
	st, err := g.calendar.Delete(in.Id)
	switch err {
	case nil:
		return &StatusResponse{
			Status: st,
			Detail: "Record has been deleted",
		}, nil
	default:
		g.logger.Error(err.Error())
		var dErr exceptions.DomainError
		if errors.As(err, &dErr) {
			return &StatusResponse{
				Status: false,
				Detail: err.Error(),
			}, nil
		} else {
			return nil, err
		}
	}
}

func (g *GRPCServer) GetEvent(ctx context.Context, in *EventIDRequest) (*EventResponse, error) {
	event, err := g.calendar.GetByID(in.Id)
	switch err {
	case nil:
		protoEvent, err := fromEven(event)
		if err != nil {
			g.logger.Error(err.Error())
			return nil, err
		}
		return &EventResponse{
			Status: true,
			Event:  protoEvent,
			Detail: "Successes",
		}, nil
	default:
		g.logger.Error(err.Error())
		var dErr exceptions.DomainError
		if errors.As(err, &dErr) {
			return &EventResponse{
				Status: false,
				Detail: err.Error(),
			}, nil
		} else {
			return nil, err
		}
	}
}

// GetEvents get several events by period
func (g *GRPCServer) GetEvents(ctx context.Context, in *PeriodRequest) (*EventsResponse, error) {
	date, err := ptypes.Timestamp(in.Date)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}

	pr := &domain.PeriodWithDate{
		Date:   date,
		Period: domain.Periods(in.Period),
	}

	events, err := g.calendar.SelectByDatePeriod(pr)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}

	results := make([]*Event, 0, len(events))
	for _, event := range events {
		protoEvent, err := fromEven(event)
		if err != nil {
			g.logger.Error(err.Error())
			return nil, err
		}
		results = append(results, protoEvent)
	}
	return &EventsResponse{
		Status: true,
		Events: results,
		Detail: "",
	}, nil
}
