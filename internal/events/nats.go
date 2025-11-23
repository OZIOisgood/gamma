package events

import (
	"log"

	"github.com/nats-io/nats.go"
)

type EventBus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewEventBus(url string) (*EventBus, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}

	log.Printf("Connected to NATS at %s (JetStream enabled)", url)
	return &EventBus{conn: nc, js: js}, nil
}

func (eb *EventBus) Publish(subject string, data []byte) error {
	_, err := eb.js.Publish(subject, data)
	return err
}

func (eb *EventBus) EnsureStream(streamName string, subjects []string) error {
	stream, err := eb.js.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		return err
	}

	if stream == nil {
		log.Printf("Creating stream %s with subjects %v", streamName, subjects)
		_, err = eb.js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (eb *EventBus) Subscribe(subject, queueGroup string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return eb.js.QueueSubscribe(subject, queueGroup, handler, nats.Durable(queueGroup), nats.ManualAck())
}

func (eb *EventBus) Close() {
	eb.conn.Close()
}
