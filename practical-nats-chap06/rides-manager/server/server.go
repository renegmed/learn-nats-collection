package server

import (
	"encoding/json"
	"log"
	"time"

	"practical-nats/rides-manager/kit"

	nats "github.com/nats-io/nats.go"
)

const (
	Version = "0.0.1"
)

type Server struct {
	*kit.Component
}

// SetupSubscriptions registers interest to the subjects that the
// Rides manager will be handling.
func (s *Server) SetupSubscription() error {
	nc := s.Component.NATS()

	// helps finding an availabe driver to accept a drive request.
	nc.QueueSubscribe("drivers.find", "manager", func(msg *nats.Msg) {
		var req *kit.DriverAgentRequest
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error on ummarshalling message data, %v", err)
			return
		}

		log.Printf("requestID:%s - Driver Find Request\n", req.RequestID)
		response := &kit.DriverAgentResponse{}

		// Find an available driver that can handle the user request
		m, err := nc.Request("drivers.rides", msg.Data, 2*time.Second)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error preparing response for driver: %s",
					req.RequestID, err)
				return
			}

			// Reply with error response
			nc.Publish(msg.Reply, resp)
			return
		}

		response.ID = string(m.Data)

		resp, err := json.Marshal(response)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error unmarshalling response from the driver: %s",
					req.RequestID, err)
				return
			}

			// Reply with error message
			nc.Publish(msg.Reply, resp)
			return
		}

		log.Printf("requestIDL%s - Driver Find Response: %+v\n",
			req.RequestID, string(m.Data))
		nc.Publish(msg.Reply, resp)

	})

	return nil
}
