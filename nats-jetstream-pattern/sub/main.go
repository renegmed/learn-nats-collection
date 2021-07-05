package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	nats "github.com/nats-io/nats.go"
)

type Message struct {
	MessageID string `json:"message-id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}

func usage() {
	log.Printf("Usage: subscriber [-s server] \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {

	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var durableName = flag.String("durablename", "subscriber", "Stream durable name")
	var streamSubject = flag.String("topic", "", "Stream subject")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	log.Println("...nats server URL", *serversUrl)

	// Connect to NATS
	nc, err := nats.Connect(*serversUrl) //nats.DefaultURL)
	if err != nil {
		log.Println("...Error on nats connection,", err)
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Println("...Error on accessing jetstream,", err)
		log.Fatal(err)
	}

	// log.Println("...Got nats jet stream connection, subscribe to", *streamSubject)

	//----------------------------- this also works

	// // Create durable consumer monitor
	// js.Subscribe(subSubjectName, func(msg *nats.Msg) {
	// 	msg.Ack()
	// 	var message Message
	// 	err := json.Unmarshal(msg.Data, &message)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	log.Printf("...Message received:\n\t%v\n", message)

	// }, nats.Durable("subscriber"), nats.ManualAck())

	// select {}

	//-----------------------------

	//+++++++++ this also works ++++++++++++
	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.

	log.Println("...Subscriber subject name:", *streamSubject)
	sub, err := js.PullSubscribe(*streamSubject, *durableName, nats.PullMaxWaiting(128))
	if err != nil {
		log.Println("...Error on pull subscribe,", err)
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Second)
	defer cancel()

	batch_size := 2

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgs, _ := sub.Fetch(batch_size, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()
			var message Message
			err := json.Unmarshal(msg.Data, &message)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("...Message received\n\tMessageID: %s\n\tMessage Topic: %s\n\tMessage: %s\n",
				message.MessageID, message.Topic, message.Message)

		}
	}

	// ++++++++++++++++++++++++++++++++++++++++++

}
