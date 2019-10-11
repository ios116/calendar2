package store

import (
	"database/sql"
	"github.com/ios116/calendar/internal/domain"
	"github.com/ios116/calendar/internal/exceptions"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

import "context"

// DbRepo - event repository
type DbRepo struct {
	ctx    context.Context
	db     *sqlx.DB
	logger *zap.Logger
}

// NewDbRepo - create event repository
func NewDbRepo(ctx context.Context, db *sqlx.DB) *DbRepo {
	return &DbRepo{ctx: ctx, db: db}
}

// AddEvent - added event to calendar list
func (d *DbRepo) AddEvent(event *domain.Event) (*domain.Event, error) {
	evDB := fromEvent(event)
	destinations := &EventDB{}
	err := d.db.GetContext(d.ctx, destinations, "INSERT INTO events (title, date, duration, author, description, notify, reminded) VALUES ($1, $2, $3, $4,$5 , $6, $7) returning *",
		evDB.Title, evDB.Date, evDB.Duration.Get(), evDB.Author, evDB.Description, evDB.Notify.Get(), evDB.Reminded)
	return toEvent(destinations), err
}

// EditEvent - edit event
func (d *DbRepo) EditEvent(event *domain.Event) (bool, error) {
	evDB := fromEvent(event)
	result, err := d.db.NamedExecContext(d.ctx, "UPDATE events SET (title, date, duration, author, description, notify, reminded) = (:title, :date, :duration, :author, :description, :notify, :reminded) WHERE id = :id",
		map[string]interface{}{
			"id":          evDB.ID,
			"title":       evDB.Title,
			"date":        evDB.Date,
			"duration":    evDB.Duration.Get(),
			"author":      evDB.Author,
			"description": evDB.Description,
			"notify":      evDB.Notify.Get(),
			"reminded":    evDB.Reminded,
		})

	if err != nil {
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, exceptions.ObjectDoesNotExist
	}

	return true, nil
}

// DeleteEvent - delete event by id
func (d *DbRepo) DeleteEvent(id int64) (bool, error) {
	res, err := d.db.ExecContext(d.ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// GetEventByID - get event by id
func (d *DbRepo) GetEventByID(id int64) (*domain.Event, error) {
	evDB := &EventDB{}
	err := d.db.GetContext(d.ctx, evDB, "SELECT * FROM events WHERE id= $1", id)

	switch err {
	case nil:
		return toEvent(evDB), nil
	case sql.ErrNoRows:
		return nil, exceptions.ObjectDoesNotExist
	default:
		return nil, err
	}
}

// SelectEventsByDatePeriod - select all event by date and period (day, week, month)
func (d *DbRepo) SelectEventsByDatePeriod(pr *domain.PeriodWithDate) ([]*domain.Event, error) {

	var events []*domain.Event
	offset, err := pr.Offset()
	if err != nil {
		return nil, err
	}
	rows, err := d.db.QueryxContext(d.ctx, "SELECT * FROM events WHERE date >= $1 AND date <= $2", pr.Date, offset)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		evDB := EventDB{}
		if err := rows.StructScan(&evDB); err != nil {
			return nil, err
		}
		events = append(events, toEvent(&evDB))
	}
	return events, nil
}

// EventReminders select events for notify
func (d *DbRepo) EventReminders(date time.Time) ([]*domain.Event, error) {
	var events []*domain.Event
	rows, err := d.db.QueryxContext(d.ctx, "SELECT * FROM events WHERE date-notify <= $1 AND reminded = FALSE ", date)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		evDB := EventDB{}
		if err := rows.StructScan(&evDB); err != nil {
			return nil, err
		}
		events = append(events, toEvent(&evDB))
	}
	return events, nil
}
