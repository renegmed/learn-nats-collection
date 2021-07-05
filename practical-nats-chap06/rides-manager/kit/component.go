package kit

import (
	"fmt"
	"sync"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

// Component contains reusable logic related
// to hanndling of the connection to NATS in the system.
type Component struct {
	// cmu is the lock form the component.
	Cmu sync.Mutex

	// id is a unique identifier used for this component.
	Id string

	// nc is the connection to NATS.
	Conn *nats.Conn

	// kind is the type of component
	Kind string
}

// NewComponent create a Component
func NewComponent(kind string) *Component {
	id := nuid.Next()
	return &Component{
		Id:   id,
		Kind: kind,
	}
}

// SetupConnectionToNATS connects to NATS
func (c *Component) SetupConnectionToNATS(servers string, options ...nats.Option) error {
	options = append(options, nats.Name(c.Name()))
	c.Cmu.Lock()
	defer c.Cmu.Unlock()

	nc, err := nats.Connect(servers, options...)
	if err != nil {
		return err
	}
	c.Conn = nc

	return nil
}

// NATS returns the current NATS connection.
func (c *Component) NATS() *nats.Conn {
	c.Cmu.Lock()
	defer c.Cmu.Unlock()
	return c.Conn
}

// ID returns the ID from the component.
func (c *Component) ID() string {
	c.Cmu.Lock()
	defer c.Cmu.Unlock()
	return c.Id
}

// Name is the label used to identify the NATS connection.
func (c *Component) Name() string {
	c.Cmu.Lock()
	defer c.Cmu.Unlock()
	return fmt.Sprintf("%s:%s", c.Kind, c.Id)
}

// Shutdown makes the component go away.
func (c *Component) Shutdown() error {
	c.NATS().Close()
	return nil
}
