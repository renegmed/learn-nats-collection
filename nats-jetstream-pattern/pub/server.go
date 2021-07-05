package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	nats "github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	MessageID string `json:"message-id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}

type server struct {
	js nats.JetStreamContext
}

func (s *server) HandlePublishMessage(rw http.ResponseWriter, req *http.Request) {

	log.Println("Request method:", req.Method)
	switch req.Method {
	case "POST":
		s.publishMessage(rw, req)
	default:
		log.Printf("Invalid reques method: %s", req.Method)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
	}
}

func (s *server) publishMessage(rw http.ResponseWriter, req *http.Request) {

	var msg Message

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("...Request Body:", string(body))

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to read request: %v", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	messageID := uuid.NewV4().String()
	msg.MessageID = messageID

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message topic %s\n\t%v", msg.Topic, err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
	_, err = s.js.Publish(msg.Topic, msgJSON)
	if err != nil {
		log.Printf("Failed to publish message onto queue '%s': %v", msg.Topic, err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
}
