package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"

	"jetstream-monitor/model"
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

	log.Println("...nats server URL", natsServers)

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...Got nats jet stream connection")

	// Create durable consumer monitor
	js.Subscribe("ORDERS.*", func(msg *nats.Msg) {
		msg.Ack()
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("...monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("...OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
	}, nats.Durable("monitor"), nats.ManualAck())

	runtime.Goexit()

}
