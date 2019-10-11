package calendar

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/ios116/calendar/internal/domain"
	"testing"
	"time"
)

func TestCalendar(t *testing.T) {
	ctrl := gomock.NewController(t)
	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	repoMock := NewMockEventRepository(ctrl)

	date := time.Now().UTC()
	duration := time.Duration(3 * time.Hour)
	notify := time.Duration(3 * time.Hour)
	event := &domain.Event{
		Title:       "Title",
		Date:        date,
		Duration:    duration,
		Author:      "Author",
		Description: "Descriptions",
		Notify:      notify,
	}

	addedEvent := event
	addedEvent.ID = 1
	var err error
	ctx := context.Background()
	calendar := NewCalendar(ctx, repoMock, nil)

	t.Run("Add", func(t *testing.T) {
		repoMock.EXPECT().AddEvent(event).Return(addedEvent, err)
		result, err := calendar.Add(event)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	})

	t.Run("Edit", func(t *testing.T) {
		repoMock.EXPECT().EditEvent(addedEvent).Return(true, err)
		result, err := calendar.Edit(addedEvent)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	})

	t.Run("Delete", func(t *testing.T) {
		repoMock.EXPECT().DeleteEvent(int64(56)).Return(true, nil)
		result, err := calendar.Delete(int64(56))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	})

}
