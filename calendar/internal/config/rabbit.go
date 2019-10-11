package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/streadway/amqp"
	"log"
)

// RabbitConf config
type RabbitConf struct {
	RQUser     string `env:"RQUser" envDefault:"guest"`
	RQPassword string `env:"RQPassword" envDefault:"123456"`
	RQHost     string `env:"RQHost" envDefault:"0.0.0.0"`
	RQPort     string `env:"RQPort" envDefault:"5672"`
}

func NewRabbitConf() *RabbitConf {
	c := &RabbitConf{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}

func RQConnection(c *RabbitConf) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", c.RQUser, c.RQPassword, c.RQHost, c.RQPort)
	var err error
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
