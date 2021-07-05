package web

import (
	"embed"
	"log"

	nats "github.com/nats-io/nats.go"
)

type Requests struct {
	Users     string
	Movies    string
	Showtimes string
	Bookings  string
}

type Application struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	Requests  Requests
	Resources *embed.FS
	Conn      *nats.Conn
}
