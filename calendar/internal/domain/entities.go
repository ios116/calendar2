package domain

import (
	"errors"
	"github.com/ios116/calendar/internal/exceptions"
	"time"
)

// Event - model for event
type Event struct {
	ID          int64
	Title       string
	Date        time.Time
	Duration    time.Duration
	Author      string
	Description string
	Notify      time.Duration
	Reminded    bool
}

// Check the correctness of the received data for the event model
func (e *Event) Validate() (err error) {
	if e.Title == "" {
		return exceptions.TitleRequired
	}
	if e.Author == "" {
		return exceptions.AuthorRequired
	}
	if e.Duration == 0 {
		return exceptions.DateRequired
	}
	if e.Date.IsZero() {
		return exceptions.DateRequired
	}
	return nil
}

// EventRepository interface for event repository
type EventRepository interface {
	AddEvent(event *Event) (*Event, error)
	EditEvent(event *Event) (bool, error)
	DeleteEvent(id int64) (bool, error)
	GetEventByID(id int64) (*Event, error)
	SelectEventsByDatePeriod(pr *PeriodWithDate) ([]*Event, error)
	EventReminders(date time.Time) ([]*Event, error)
}

// Periods - periods for a selects day, week, month
type Periods int32

// enum for periods
const (
	PeriodDay Periods = iota
	PeriodWeek
	PeriodMonth
)

// PeriodWithDate struct for filter by date and period
type PeriodWithDate struct {
	Date   time.Time
	Period Periods
}

// Offset is adding a period to passed date
func (p *PeriodWithDate) Offset() (time.Time, error) {
	var date time.Time
	switch p.Period {
	case PeriodDay:
		date = p.Date.AddDate(0, 0, 1)
	case PeriodWeek:
		date = p.Date.AddDate(0, 0, 7)
	case PeriodMonth:
		date = p.Date.AddDate(0, 1, 0)
	default:
		return date, errors.New("period isn't validate")
	}
	return date, nil
}
