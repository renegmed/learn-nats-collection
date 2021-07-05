package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"practical-nats/riders-client/kit"

	"github.com/nats-io/nuid"
)

const (
	Version = "0.0.1"
)

type Server struct {
	*kit.Component
}

func NewServer(c *kit.Component) Server {
	return Server{Component: c}
}

func (s *Server) HandleRequestRides(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		s.requestDriver(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) requestDriver(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request *kit.DriverAgentRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Tag the request with an ID for tacing in the logs
	request.RequestID = nuid.Next()
	req, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	nc := s.Component.NATS()

	log.Printf("requestID:%s - Finding available driver for request: %s\n", request.RequestID, string(body))

	msg, err := nc.Request("drivers.find", req, 5*time.Second)
	if err != nil {
		log.Printf("requestID:%s - Gave up finding available driver for request\n", request.RequestID)
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	log.Printf("requestID:%s - Response: %s\n", request.RequestID, string(msg.Data))

	var resp *kit.DriverAgentResponse
	err = json.Unmarshal(msg.Data, &resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if resp.Error != "" {
		http.Error(w, resp.Error, http.StatusServiceUnavailable)
		return
	}

	log.Printf("requestID:%s - Driver with ID %s is available to handle the request", request.RequestID, resp.ID)

	fmt.Fprint(w, string(msg.Data))
}

func (s *Server) ListenAndServe(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, fmt.Sprintf("NATS Rider API Server v%s\n", Version))
	})

	mux.HandleFunc("/rides", s.HandleRequestRides)

	addr := "localhost:" + port

	log.Println("...Server running on localhost:", port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error on  initializing port listener, %v", err)
	}

	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go srv.Serve(l)

	return nil
}
