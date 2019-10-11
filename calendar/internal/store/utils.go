package store

import (
	"database/sql"
	"github.com/ios116/calendar/internal/domain"
	"github.com/jackc/pgx/pgtype"
	"time"
)

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toEvent(evDB *EventDB) *domain.Event {
	var description string
	if evDB.Description.Valid {
		description = ""
	}
	return &domain.Event{
		ID:          evDB.ID,
		Title:       evDB.Title,
		Date:        evDB.Date,
		Duration:    toInterval(evDB.Duration),
		Author:      evDB.Author,
		Description: description,
		Notify:      toInterval(evDB.Notify),
		Reminded:    evDB.Reminded,
	}
}

func fromEvent(ev *domain.Event) *EventDB {
	var duration pgtype.Interval
	if err := duration.Set(ev.Duration); err != nil {

	}
	var notify pgtype.Interval
	if err := notify.Set(ev.Notify); err != nil {

	}
	return &EventDB{
		ID:          ev.ID,
		Title:       ev.Title,
		Date:        ev.Date,
		Duration:    duration,
		Author:      ev.Author,
		Description: toNullString(ev.Description),
		Notify:      notify,
		Reminded:    ev.Reminded,
	}
}

const (
	day   int64 = int64(time.Hour) * int64(24)
	month int64 = day * 30
)

func toInterval(t pgtype.Interval) time.Duration {
	dur := time.Duration(t.Microseconds*int64(time.Microsecond) + int64(t.Days)*day + int64(t.Months)*month)
	return dur
}
