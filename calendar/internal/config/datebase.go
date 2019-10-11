package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jmoiron/sqlx"
	"log"
)

// DateBaseConf config
type DateBaseConf struct {
	BdPassword string `env:"PGPASSWORD" envDefault:"123456"`
	BdUser     string `env:"PGUSER" envDefault:"calendar"`
	BdHost     string `env:"PGHOST" envDefault:"0.0.0.0"`
	BdName     string `env:"PGDATABASE" envDefault:"calendar"`
}

func NewDateBaseConf() *DateBaseConf {
	c := &DateBaseConf{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}

// DBConnection - connection for BD
// postgres://myuser:mypass@localhost:5432/mydb?sslmode=verify­full
// export POSTGRESQL_URL=postgres://calendar:123456@localhost:5432/calendar?sslmode=disable
func DBConnection(c *DateBaseConf) (*sqlx.DB, error) {
	var params = fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable", c.BdUser, c.BdName, c.BdHost, c.BdPassword)
	db, err := sqlx.Connect("pgx", params)
	return db, err
}
