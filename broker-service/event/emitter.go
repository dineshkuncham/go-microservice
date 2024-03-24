package event

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	c, err := e.connection.Channel()
	if err != nil {
		return err
	}
	return declareExhange(c)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	log.Println("Pushing to channel")
	err = channel.PublishWithContext(context.TODO(),
		"logs_topic", severity, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		})
	if err != nil {
		return err
	}
	return nil
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emiter := Emitter{
		connection: conn,
	}
	err := emiter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emiter, nil
}
