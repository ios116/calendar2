package store

import (
	"database/sql"
	"github.com/jackc/pgx/pgtype"
	"time"
)

// EventDB - event model for postgres
type EventDB struct {
	ID          int64
	Title       string
	Date        time.Time
	Duration    pgtype.Interval
	Author      string
	Description sql.NullString
	Notify      pgtype.Interval
	Reminded    bool
}
