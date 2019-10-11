package sender

import (
	"encoding/json"
	"github.com/ios116/calendar/internal/calendar"
	"github.com/ios116/calendar/internal/config"
	"github.com/ios116/calendar/internal/domain"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Consumer - interface
type Mailer interface {
	Send(msg interface{}) error
}

// Sender - mailer object
type Sender struct {
	mail    Mailer
	Logger  *zap.Logger
	conn    *amqp.Connection
	service calendar.UseCaseCalendar
}

// NewSender - constructor
func NewSender(mail Mailer, logger *zap.Logger, conn *amqp.Connection, service calendar.UseCaseCalendar) *Sender {
	return &Sender{mail: mail, Logger: logger, conn: conn, service: service}
}

// Consume - implementation queue
func (m *Sender) Consume() {

	sugar := m.Logger.Sugar()
	conf := config.NewRabbitConf()
	conn, err := config.RQConnection(conf)

	if err != nil {
		sugar.Fatal(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		sugar.Fatal(err)
	}
	defer ch.Close()

	sugar.Info("Connection is open...")
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		sugar.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,    // queue name
		"sending", // routing key
		"remind",  // exchange
		false,
		nil)
	if err != nil {
		sugar.Fatal(err)
	}
	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		sugar.Fatal(err)
	}
	for d := range messages {
		event := &domain.Event{}
		err = json.Unmarshal(d.Body, event)
		if err != nil {
			sugar.Error(err)
			continue
		}
		event.Reminded = true
		if err := m.mail.Send(event); err != nil {
			sugar.Error(err)
		}
		if _, err := m.service.Edit(event); err != nil {
			sugar.Error(err)
		}
	}
}
