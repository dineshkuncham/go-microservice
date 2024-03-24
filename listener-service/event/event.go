package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExhange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // kind
		true,         //durable
		false,        //auto delete
		false,        //internal?
		false,        //no-wait?
		nil,          //args
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}
