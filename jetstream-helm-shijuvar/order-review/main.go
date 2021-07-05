package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"

	"jsdemo-shijuvar/model"
)

const (
	subSubjectName = "ORDERS.created"
	pubSubjectName = "ORDERS.approved"
)

func main() {
	var (
		natsServers string
	)

	flag.StringVar(&natsServers, "s", nats.DefaultURL, "List of NATS Servers to connect")
	flag.StringVar(&natsServers, "servers", nats.DefaultURL, "List of NATS Servers to connect")
	flag.Parse()

	log.Println("...nats servers URL", natsServers)

	// Connect to NATS
	nc, err := nats.Connect(natsServers)
	if err != nil {
		log.Println("...Error on nats connection,", err)
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...Got nats jet stream connection")

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	sub, err := js.PullSubscribe(subSubjectName, "order-review", nats.PullMaxWaiting(128))
	if err != nil {
		log.Println("...Error on pull subscribe,", err)
		log.Fatal(err)
	}

	log.Println("...Pull subscribe with subscription subject name:", subSubjectName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgs, _ := sub.Fetch(10, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()
			var order model.Order
			err := json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("...order-review service")
			log.Printf("...OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
			reviewOrder(js, order)
		}
	}
}

// reviewOrder reviews the order and publishes ORDERS.approved event
func reviewOrder(js nats.JetStreamContext, order model.Order) {
	// Changing the Order status
	order.Status = "approved"
	orderJSON, _ := json.Marshal(order)
	_, err := js.Publish(pubSubjectName, orderJSON)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("...Publish on subject '%s' Order with OrderID:%d has been %s\n", pubSubjectName, order.OrderID, order.Status)
}
