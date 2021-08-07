package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"runtime"

	"cqrs-drones/events/common"
	"cqrs-drones/events/mongo"

	nats "github.com/nats-io/nats.go"
)

type drone struct {
	ID string `json:"id"`
}

func execute() error {

	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var dbUrl = flag.String("dburl", nats.DefaultURL, "Database URL address")

	flag.Parse()

	// Connect to NATS
	opts := []nats.Option{nats.Name("NATS Drone event generation")}
	conn, err := nats.Connect(*serversUrl, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dbase, err := mongo.Connect(*dbUrl)
	if err != nil {
		return fmt.Errorf("Problem setting up connection to database, %v", err)
	}
	repo := mongo.NewEventRollupRepository(dbase)

	go subscribe(conn, repo)

	select {}

	return nil
}

func subscribe(conn *nats.Conn, repo *mongo.EventRollupRepository) {

	conn.Subscribe("New.Alert.Signal", func(m *nats.Msg) {
		var event common.AlertSignalledEvent
		err := json.Unmarshal(m.Data, &event)
		if err != nil {
			log.Println("Error on unmarshalling alert signal event", err)
			return
		}

		log.Printf("Received alert signal update event,\n\t%v", event)
		err = repo.UpdateLastAlertEvent(event)
		if err != nil {
			log.Println("Error on updating repository on alert signal event,", err)
			return
		}
		m.Respond([]byte("New.Alert.Signal data received."))
	})

	conn.Subscribe("Drone.Alert.Signal", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id for alert signal,", err)
			return
		}

		log.Printf("Received request drone alert signal,\n\t%v", drone)

		alertSignalledEvent, err := repo.GetAlertEvent(drone.ID)
		if err != nil {
			log.Println("Error on getting drone alert signal event from repository,", err)
			return
		}

		log.Printf("Drone alert signal event:\n\t%v", alertSignalledEvent)

		alertSignal, err := json.Marshal(alertSignalledEvent)
		if err != nil {
			log.Println("Error on marshalling alert signal data,", err)
			return
		}
		m.Respond(alertSignal)
	})

	conn.Subscribe("New.Telemetry", func(m *nats.Msg) {
		var event common.TelemetryUpdatedEvent
		err := json.Unmarshal(m.Data, &event)
		if err != nil {
			log.Println("Error on unmarshalling telemetry update event", err)
			return
		}

		log.Printf("Received telemetry update event,\n\t%v", event)
		err = repo.UpdateLastTelemetryEvent(event)
		if err != nil {
			log.Println("Error on updating repository on telemetry update event,", err)
			return
		}
		m.Respond([]byte("New.Telemetry data received."))
	})

	conn.Subscribe("Drone.Telemetry", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id for telemetry,", err)
			return
		}

		log.Printf("Received request drone telemetry,\n\t%v", drone)

		telemetryEvent, err := repo.GetTelemetryEvent(drone.ID)
		if err != nil {
			log.Println("Error on getting telemetry from repository,", err)
			return
		}

		log.Printf("Drone telemetry:\n\t%v", telemetryEvent)

		droneTelemetry, err := json.Marshal(telemetryEvent)
		if err != nil {
			log.Println("Error on marshalling telemetry data,", err)
			return
		}
		m.Respond(droneTelemetry)
	})

	conn.Subscribe("Change.Position", func(m *nats.Msg) {
		var event common.PositionChangedEvent
		err := json.Unmarshal(m.Data, &event)
		if err != nil {
			log.Println("Error on unmarshalling position changed event", err)
			return
		}

		log.Printf("Received changed position event,\n\t%v", event)
		err = repo.UpdateLastPositionEvent(event)
		if err != nil {
			log.Println("Error on updating repository on changed position event,", err)
			return
		}
		m.Respond([]byte("Change.Position data received."))
	})

	conn.Subscribe("Last.Position", func(m *nats.Msg) {
		var drone drone
		err := json.Unmarshal(m.Data, &drone)
		if err != nil {
			log.Println("Error on unmarshalling drone id,", err)
			return
		}

		log.Printf("Received request drone last position,\n\t%v", drone)

		lastPosition, err := repo.GetPositionEvent(drone.ID)
		if err != nil {
			log.Println("Error on getting drone last position,", err)
			return
		}

		log.Printf("Drone last position:\n\t%v", lastPosition)

		dronePosition, err := json.Marshal(lastPosition)
		if err != nil {
			log.Println("Error on marshalling drone position,", err)
			return
		}
		m.Respond(dronePosition)
	})
}

func main() {

	err := execute()
	if err != nil {
		log.Println(err)
	}
	runtime.Goexit()
}
