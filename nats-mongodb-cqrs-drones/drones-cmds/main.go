package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	nats "github.com/nats-io/nats.go"
)

type drone struct {
	ID string `json:"id"`
}

type response struct {
	DroneID string `json:"droneid"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type server struct {
	nc *nats.Conn
}

var natsServer server

func addTelemetryHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		log.Println("Requesting to add telemetry to queue")

		payload, _ := ioutil.ReadAll(req.Body)
		var newTelemetryCommand telemetryCommand
		err := json.Unmarshal(payload, &newTelemetryCommand)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: "000", Type: "fail", Message: "Failed to parse add telemetry command."})
			return
		}
		if !newTelemetryCommand.isValid() {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newTelemetryCommand.DroneID, Type: "fail", Message: "Invalid add telemetry command."})
			return
		}

		evt := TelemetryUpdatedEvent{
			DroneID:          newTelemetryCommand.DroneID,
			RemainingBattery: newTelemetryCommand.RemainingBattery,
			Uptime:           newTelemetryCommand.Uptime,
			CoreTemp:         newTelemetryCommand.CoreTemp,
			ReceivedOn:       time.Now().Unix(),
		}

		evtBody, err := json.Marshal(evt)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newTelemetryCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not marshal telemetry data, %v", err)})
			return
		}

		fmt.Printf("Dispatching telemetry event for drone %s\n", newTelemetryCommand.DroneID)

		_, err = natsServer.nc.Request("New.Telemetry", evtBody, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newTelemetryCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not add new telemetry, %v", err)})
			return
		}
		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: newTelemetryCommand.DroneID, Type: "success", Message: "telemetry added."})
	}
}

func queryTelemetryHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("query position handler req.Method,", req.Method)

	switch req.Method {
	case http.MethodGet:

		urlPathSegments := strings.Split(req.URL.Path, "query-telemetry/")
		id := urlPathSegments[len(urlPathSegments)-1] // get the last part of array

		log.Printf("Drone ID query request, %v", id)

		drone := drone{id}

		payload, err := json.Marshal(drone)
		if err != nil {
			serviceResponse(w, http.StatusInternalServerError,
				response{DroneID: id, Type: "fail", Message: fmt.Sprintf("Could not marshal drone request, %v", err)})
			return
		}
		resp, err := natsServer.nc.Request("Query.Telemetry", payload, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: drone.ID, Type: "fail", Message: fmt.Sprintf("Could not query drone telemetry, %v", err)})
			return
		}

		log.Printf("Drone query response, %", string(resp.Data))

		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: drone.ID, Type: "success", Message: string(resp.Data)})
	}
}

func changedPositionHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		log.Println("Requesting to change position to queue")

		payload, _ := ioutil.ReadAll(req.Body)
		var newPositionCommand positionCommand
		err := json.Unmarshal(payload, &newPositionCommand)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: "000", Type: "fail", Message: "Failed to parse changed position command."})
			return
		}
		if !newPositionCommand.isValid() {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: "000", Type: "fail", Message: "Invalid change position command."})
			return
		}
		evt := PositionChangedEvent{
			DroneID:         newPositionCommand.DroneID,
			Longitude:       newPositionCommand.Longitude,
			Latitude:        newPositionCommand.Latitude,
			Altitude:        newPositionCommand.Altitude,
			CurrentSpeed:    newPositionCommand.CurrentSpeed,
			HeadingCardinal: newPositionCommand.HeadingCardinal,
			ReceivedOn:      time.Now().Unix(),
		}

		evtBody, err := json.Marshal(evt)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newPositionCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not marshal data, %v", err)})
			return
		}

		fmt.Printf("Dispatching change position event for drone %s\n", newPositionCommand.DroneID)

		_, err = natsServer.nc.Request("Change.Position", evtBody, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newPositionCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not add changed position data, %v", err)})
			return
		}
		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: newPositionCommand.DroneID, Type: "success", Message: "changed position data added."})
	}
}

func queryPositionHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("query position hadler req.Method,", req.Method)

	switch req.Method {
	case http.MethodGet:

		urlPathSegments := strings.Split(req.URL.Path, "query-position/")
		id := urlPathSegments[len(urlPathSegments)-1] // get the last part of array

		log.Printf("Drone ID query request, %v", id)

		drone := drone{id}

		payload, err := json.Marshal(drone)
		if err != nil {
			serviceResponse(w, http.StatusInternalServerError,
				response{DroneID: id, Type: "fail", Message: fmt.Sprintf("Could not marshal drone request, %v", err)})
			return
		}
		resp, err := natsServer.nc.Request("Query.Last.Position", payload, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: drone.ID, Type: "fail", Message: fmt.Sprintf("Could not query drone last position, %v", err)})
			return
		}

		log.Printf("Drone query response, %", string(resp.Data))

		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: drone.ID, Type: "success", Message: string(resp.Data)})
	}
}

func addAlertSignalHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		log.Println("Requesting store alert signal")

		payload, _ := ioutil.ReadAll(req.Body)

		var newAlertCommand alertCommand
		err := json.Unmarshal(payload, &newAlertCommand)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: "000", Type: "fail", Message: "Failed to parse add alert command."})
			return
		}
		if !newAlertCommand.isValid() {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: "000", Type: "fail", Message: "Invalid alert command."})
			return
		}
		evt := AlertSignalledEvent{
			DroneID:     newAlertCommand.DroneID,
			FaultCode:   newAlertCommand.FaultCode,
			Description: newAlertCommand.Description,
			ReceivedOn:  time.Now().Unix(),
		}

		evtBody, err := json.Marshal(evt)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newAlertCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not marshal data, %v", err)})
			return
		}

		fmt.Printf("Dispatching alert signal event for drone %s\n", newAlertCommand.DroneID)

		_, err = natsServer.nc.Request("New.Alert.Signal", evtBody, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: newAlertCommand.DroneID, Type: "fail", Message: fmt.Sprintf("Could not store alert signal data, %v", err)})
			return
		}
		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: newAlertCommand.DroneID, Type: "success", Message: "alert signal data added."})
	}
}

func queryAlertSignalHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("query position handler req.Method,", req.Method)

	switch req.Method {
	case http.MethodGet:

		urlPathSegments := strings.Split(req.URL.Path, "query-alertsignal/")
		id := urlPathSegments[len(urlPathSegments)-1] // get the last part of array

		log.Printf("Drone ID query request, %v", id)

		drone := drone{id}

		payload, err := json.Marshal(drone)
		if err != nil {
			serviceResponse(w, http.StatusInternalServerError,
				response{DroneID: id, Type: "fail", Message: fmt.Sprintf("Could not marshal drone request, %v", err)})
			return
		}
		resp, err := natsServer.nc.Request("Query.Alert.Signal", payload, 1*time.Second)
		if err != nil {
			serviceResponse(w, http.StatusBadRequest,
				response{DroneID: drone.ID, Type: "fail", Message: fmt.Sprintf("Could not query drone alert signal, %v", err)})
			return
		}

		log.Printf("Drone query response, %", string(resp.Data))

		serviceResponse(w, http.StatusBadRequest,
			response{DroneID: drone.ID, Type: "success", Message: string(resp.Data)})
	}
}

func main() {

	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var serverPort = flag.String("port", "8080", "This server port number e.g. 8080")

	flag.Parse()

	// Connect to NATS
	opts := []nats.Option{nats.Name("NATS Drone Query")}
	conn, err := nats.Connect(*serversUrl, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	natsServer = server{nc: conn}

	http.HandleFunc("/api/cmds/telemetry", addTelemetryHandler)
	http.HandleFunc("/api/cmds/query-telemetry/", queryTelemetryHandler)
	http.HandleFunc("/api/cmds/position", changedPositionHandler)
	http.HandleFunc("/api/cmds/query-position/", queryPositionHandler)
	http.HandleFunc("/api/cmds/alertsignal", addAlertSignalHandler)
	http.HandleFunc("/api/cmds/query-alertsignal/", queryAlertSignalHandler)

	log.Printf("====== Cmds server listening on port %s...", *serverPort)
	if err := http.ListenAndServe(":"+*serverPort, nil); err != nil {
		log.Fatal(err)
	}
}

func serviceResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(body)
}
