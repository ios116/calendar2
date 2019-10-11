package scheduler

import (
	"encoding/json"
	"github.com/ios116/calendar/internal/calendar"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

// Scanner - scanner object
type Scanner struct {
	service calendar.UseCaseCalendar
	conn    *amqp.Connection
	logger  *zap.Logger
}

// NewScanner - constructor
func NewScanner(service calendar.UseCaseCalendar, conn *amqp.Connection, logger *zap.Logger) *Scanner {
	return &Scanner{service: service, conn: conn, logger: logger}
}

// Produce - implementation queue
func (m *Scanner) Produce() {
	sugar := m.logger.Sugar()
	if m.conn == nil {
		sugar.Fatal("Connection is nil")
	}
	defer m.conn.Close()
	ch, err := m.conn.Channel()
	if err != nil {
		sugar.Fatal(err)
	}
	defer ch.Close()

	sugar.Info("Scanner is start")

	err = ch.ExchangeDeclare(
		"remind", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		sugar.Fatal(err)
	}

	for {
		date := <-time.Tick(5 * time.Second)
		events, err := m.service.EventReminders(date)
		if err != nil {
			sugar.Fatal(err)
		}

		sugar.Info("scanning", events)
		for _, event := range events {
			payload, err := json.Marshal(event)
			err = ch.Publish(
				"remind",  // exchange
				"sending", // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        payload,
				})
			if err != nil {
				sugar.Fatal(err)
			}
		}
	}
}
