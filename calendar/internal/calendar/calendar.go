package calendar

import (
	"context"
	"github.com/ios116/calendar/internal/domain"
	"github.com/ios116/calendar/internal/exceptions"
	"go.uber.org/zap"
	"time"
)

// Calendar - service
type Calendar struct {
	ctx        context.Context
	repository domain.EventRepository
	logger     *zap.Logger
}

// NewCalendar - constructor for Calendar
func NewCalendar(ctx context.Context, repository domain.EventRepository, logger *zap.Logger) *Calendar {
	return &Calendar{repository: repository, logger: logger, ctx: ctx}
}

// UseCaseCalendar - use case for calendar
type UseCaseCalendar interface {
	Add(event *domain.Event) (*domain.Event, error)
	Edit(event *domain.Event) (bool, error)
	Delete(id int64) (bool, error)
	GetByID(id int64) (*domain.Event, error)
	SelectByDatePeriod(pr *domain.PeriodWithDate) ([]*domain.Event, error)
	EventReminders(date time.Time) ([]*domain.Event, error)
}

// Add - added event to calendar list
func (c *Calendar) Add(event *domain.Event) (*domain.Event, error) {
	if err := event.Validate(); err != nil {
		return nil, err
	}
	return c.repository.AddEvent(event)
}

// Edit - edit event
func (c *Calendar) Edit(event *domain.Event) (bool, error) {
	if err := event.Validate(); err != nil {
		return false, err
	}
	return c.repository.EditEvent(event)
}

// Delete - delete event by id
func (c *Calendar) Delete(id int64) (bool, error) {
	if id == 0 {
		return false, exceptions.IDRequired
	}
	return c.repository.DeleteEvent(id)
}

// GetByID - get event by id
func (c *Calendar) GetByID(id int64) (*domain.Event, error) {
	if id == 0 {
		return nil, exceptions.IDRequired
	}
	return c.repository.GetEventByID(id)
}

// SelectByDatePeriod - select all event by date and period (day, week, month)
func (c *Calendar) SelectByDatePeriod(pr *domain.PeriodWithDate) ([]*domain.Event, error) {
	return c.repository.SelectEventsByDatePeriod(pr)
}
// EventReminders - select all events for remindt brnch

func (c *Calendar) EventReminders(date time.Time) ([]*domain.Event, error) {
	if date.IsZero() {
		return nil, exceptions.DateRequired
	}
	return c.repository.EventReminders(date)
}
