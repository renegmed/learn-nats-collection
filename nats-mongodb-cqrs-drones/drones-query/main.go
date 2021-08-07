package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	nats "github.com/nats-io/nats.go"
)

type drone struct {
	ID string `json:"id"`
}

var usageStr = `
Usage: product-read [options] <subject> <message>

Options:  
	-s,            NATS stream server URL(s)
 
`

func usage() {
	log.Printf("%s\n", usageStr)
	os.Exit(0)
}

func execute() error {

	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

	flag.Usage = usage
	flag.Parse()

	// Connect to NATS
	opts := []nats.Option{nats.Name("NATS Drone Query")}
	conn, err := nats.Connect(*serversUrl, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go subscribe(conn)

	select {}

	return nil
}

func subscribe(conn *nats.Conn) {

	conn.Subscribe("Query.Last.Position", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id", err)
			return
		}

		data, err := json.Marshal(drone)
		if err != nil {
			log.Println("Error on marshalling drone request data,", err)
			return
		}

		log.Printf("Received query last position of drone,\n\t%v", drone)

		resp, err := conn.Request("Last.Position", data, 500*time.Millisecond)
		if err != nil {
			log.Println("Error on request last position,", err)
			return
		}
		log.Printf("Response Data: \n\t%v", string(resp.Data))

		m.Respond(resp.Data)
	})

	conn.Subscribe("Query.Telemetry", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id", err)
			return
		}

		data, err := json.Marshal(drone)
		if err != nil {
			log.Println("Error on marshalling drone request data,", err)
			return
		}

		log.Printf("Received query drone's telemetry data ,\n\t%v", drone)

		resp, err := conn.Request("Drone.Telemetry", data, 500*time.Millisecond)
		if err != nil {
			log.Println("Error on request drone's telemetry data,", err)
			return
		}
		log.Printf("Response Data: \n\t%v", string(resp.Data))

		m.Respond(resp.Data)
	})

	conn.Subscribe("Query.Alert.Signal", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id", err)
			return
		}

		data, err := json.Marshal(drone)
		if err != nil {
			log.Println("Error on marshalling drone request data,", err)
			return
		}

		log.Printf("Received query drone's alert signal ,\n\t%v", drone)

		resp, err := conn.Request("Drone.Alert.Signal", data, 500*time.Millisecond)
		if err != nil {
			log.Println("Error on request drone's alert signal,", err)
			return
		}
		log.Printf("Response Data: \n\t%v", string(resp.Data))

		m.Respond(resp.Data)
	})
}

func main() {

	err := execute()
	if err != nil {
		log.Println(err)
	}
	runtime.Goexit()
}
