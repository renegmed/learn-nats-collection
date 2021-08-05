package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	nats "github.com/nats-io/nats.go"

<<<<<<< HEAD
	"ordering-app/model"
=======
	"jetstream-order/model"
>>>>>>> c39dd6d... added crud for showtimes and other modifications
)

const (
	streamName     = "ORDERS"
	streamSubjects = "ORDERS.*"
	subjectName    = "ORDERS.created"
)

func main() {
	var (
		natsServers string
	)

	flag.StringVar(&natsServers, "s", nats.DefaultURL, "List of NATS Servers to connect")
	flag.StringVar(&natsServers, "servers", nats.DefaultURL, "List of NATS Servers to connect")
	flag.Parse()
	// Connect to NATS
	nc, err := nats.Connect(natsServers)
	if err != nil {
		log.Println("...Error on nats connection,", err)
		log.Fatal(err)
	}

	log.Println("...nats servers URL", natsServers)

	// Creates JetStreamContext
	js, err := nc.JetStream()
	checkErr(err)

	log.Println("...Got nats jet stream connection")

	// Creates stream
	err = createStream(js)
	checkErr(err)

	log.Println("...stream created with stream name -", streamName, " and subjects:\n\t", streamSubjects)

	// Create orders by publishing messages
	err = createOrder(js)
	checkErr(err)
}

// createOrder publishes stream of events
// with subject "ORDERS.created"
func createOrder(js nats.JetStreamContext) error {
	var order model.Order
	//for i := 1; i <= 10; i++ {
	for {
		i := rand.Intn(1000) + 100

		order = model.Order{
			OrderID:    i,
			CustomerID: "Cust-" + strconv.Itoa(i),
			Status:     "created",
		}
		orderJSON, _ := json.Marshal(order)
		_, err := js.Publish(subjectName, orderJSON)
		if err != nil {
			return err
		}
		log.Printf("Order with OrderID:%d has been published\n", i)

		s := rand.Intn(5000) + 200

		time.Sleep(time.Duration(s) * time.Millisecond)

	}
	return nil
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
