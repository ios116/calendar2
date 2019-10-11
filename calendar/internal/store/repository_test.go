package store

import (
	"context"
	"github.com/ios116/calendar/internal/config"
	"github.com/ios116/calendar/internal/domain"
	"github.com/ios116/calendar/internal/exceptions"
	"log"
	"os"
	"testing"
	"time"
)

var repo *DbRepo

func TestMain(m *testing.M) {
	conf := config.NewDateBaseConf()
	db, err := config.DBConnection(conf)
	if err != nil {
		log.Fatal("create store with err", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("create store with err", err)
	}
	repo = NewDbRepo(context.Background(), db)
	os.Exit(m.Run())
}

func TestCrud(t *testing.T) {
	date := time.Now().UTC()
	event := &domain.Event{
		Title:       "title",
		Date:        date,
		Duration:    time.Duration(3 * 24 * time.Hour),
		Author:      "Vladimir",
		Description: "",
		Notify:      time.Duration(24 * time.Hour),
	}
	t.Run("insert", func(t *testing.T) {
		newEvent, err := repo.AddEvent(event)
		if err != nil {
			t.Fatal(err)
		}
		if newEvent.Title != "title" {
			t.Fatal("title is not equal")
		}
		event.ID = newEvent.ID

	})
	t.Run("note exist", func(t *testing.T) {
		_, err := repo.GetEventByID(5558)
		if err != exceptions.ObjectDoesNotExist {
			t.Fatal(err)
		}
	})
	t.Run("Edit", func(t *testing.T) {
		_, err := repo.EditEvent(event)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		_, err := repo.DeleteEvent(event.ID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DELETE IF NOT EXIST", func(t *testing.T) {
		_, err := repo.DeleteEvent(666)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("select by periods", func(t *testing.T) {
		oldDate := date.Add(-time.Duration(20 * time.Hour))
		pr := &domain.PeriodWithDate{
			Date:   oldDate,
			Period: domain.PeriodWeek,
		}
		if res, err := repo.SelectEventsByDatePeriod(pr); err != nil {
			t.Fatal(err)
		} else {
			t.Log(res)
		}
	})

	t.Run("events for reminded", func(t *testing.T) {
		event.Notify = time.Duration(24 * time.Hour)
		event.Date = time.Now().Add(time.Duration(2 * 24 * time.Hour))
		repo.AddEvent(event)

		event.Notify = time.Duration(2 * 24 * time.Hour)
		event.Date = time.Now().Add(time.Duration(24 * time.Hour))
		repo.AddEvent(event)

		date := time.Now()
		if res, err := repo.EventReminders(date); err != nil {
			t.Fatal(err)
		} else {
			t.Log(res)
		}
	})
}
