package test

import (
	"fmt"
	"log"
	"time"

	"practical-nats/riders-client/kit"

	stan "github.com/nats-io/stan.go"
)

func Publish(c *kit.Component, id string) {
	for i := 0; i < 5; i++ {

		msg := fmt.Sprintf("%s msg %d", id, i)
		if err := c.Connection().Publish("goo", []byte(msg)); err != nil {
			log.Println("failed to publish '", msg, "':", err)
			continue
		}
		i++
		log.Print("sending:", msg)

	}

	_, err := c.Connection().QueueSubscribe("goo", "gfoo", func(m *stan.Msg) {
		if err := m.Ack(); err != nil {
			log.Println(err)
		}
		// fake processing time
		time.Sleep(time.Millisecond * 10)
		log.Println("received:", string(m.Data))
	}, stan.MaxInflight(10), stan.AckWait(time.Second), stan.SetManualAckMode())
	if err != nil {
		log.Fatalln(err)
	}
	//defer sub.Unsubscribe()
}
