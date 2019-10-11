package grpcserver

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/ios116/calendar/internal/domain"
)

func toEvent(in *Event) (*domain.Event, error) {
	date, err := ptypes.Timestamp(in.Date)
	if err != nil {
		return nil, err
	}

	duration, err := ptypes.Duration(in.Duration)
	if err != nil {
		return nil, err
	}
	notify, err := ptypes.Duration(in.Notify)
	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:          in.Id,
		Title:       in.Title,
		Date:        date,
		Duration:    duration,
		Author:      in.Author,
		Description: in.Description,
		Notify:      notify,
		Reminded:    in.Reminded,
	}, nil
}

func fromEven(in *domain.Event) (*Event, error) {

	date, err := ptypes.TimestampProto(in.Date)
	if err != nil {
		return nil, err
	}

	return &Event{
		Id:          in.ID,
		Title:       in.Title,
		Date:        date,
		Duration:    ptypes.DurationProto(in.Duration),
		Author:      in.Author,
		Description: in.Description,
		Notify:      ptypes.DurationProto(in.Notify),
		Reminded:    in.Reminded,
	}, nil
}
