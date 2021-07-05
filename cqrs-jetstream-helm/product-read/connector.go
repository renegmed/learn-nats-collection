package main

import (
	"log"

	nats "github.com/nats-io/nats.go"
)

type Connector struct {
	Js nats.JetStreamContext
}

func NewConnector(URL string) (*Connector, error) {

	// Connect to NATS
	nc, err := nats.Connect(URL)
	if err != nil {
		log.Println("...Error on nats connection,", err)
		return &Connector{}, err
	}

	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		log.Println("...Error on accessing jetstream,", err)
		return &Connector{}, err
	}

	return &Connector{
		Js: js,
	}, nil
}
