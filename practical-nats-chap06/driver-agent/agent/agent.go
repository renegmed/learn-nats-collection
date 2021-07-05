package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"practical-nats/driver-agent/kit"

	nats "github.com/nats-io/nats.go"
)

const Version = "0.0.1"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Agent is the agent from the driver that provides rides.
type Agent struct {
	*kit.Component

	// mu is the lock from the agent
	mu sync.Mutex

	// AgentType is the type of vehicle
	AgentType string
}

// Type returns the type of vehicle form the driver
func (s *Agent) Type() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.AgentType
}

// SetupSubscriptions prepares the NATS subscriptions.
func (s *Agent) SetupSubscriptions() error {
	nc := s.Component.NATS()

	nc.Subscribe("drivers.rides", func(msg *nats.Msg) {
		if err := s.processRequest(msg); err != nil {
			log.Printf("Error: %s\n", err)
			return
		}
	})

	return nil
}

func (s *Agent) processRequest(msg *nats.Msg) error {
	var req *kit.DriverAgentRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Printf("Error while unmarshalling request data:%v\n", err)
		return fmt.Errorf("Error while unmarshalling request data:%v\n", err)
	}

	log.Printf("requestID:%s - Driver Ride Request: type: %s\n",
		req.RequestID, req.Type)

	log.Print("Agent type:", s.Type())

	if req.Type != s.Type() {
		log.Printf("Skip request since agent is of different type, %s\n", s.Type())
		return nil
	}

	log.Printf("requestID:%s - Avaible to handle request", req.RequestID)

	// Randomly delay agent when receiving drive request
	// to siumlatelatency in replying.
	duration := time.Duration(rand.Int31n(1000)) * time.Millisecond
	log.Print("requestID:%s - Backing off for %s before replying", req.RequestID, duration)
	time.Sleep(duration)

	// Replying back with own ID meaning that can help.
	return s.Component.NATS().Publish(msg.Reply, []byte(s.Component.Id))
}
