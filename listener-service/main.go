package main

import (
	"log"
	"math"
	"time"

	"github.com/dineshkuncham/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connection, err := connectToRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer connection.Close()
	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(connection)
	if err != nil {
		log.Panic(err)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

func connectToRabbitMQ() (*amqp.Connection, error) {
	connection := &amqp.Connection{}
	counts := 0
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ is not ready yet...")
			counts += 1
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if counts > 5 {
			log.Println("Unable to connect rabbitMQ", err)
			return nil, err
		}
		incrementalBackOff := time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(incrementalBackOff)
		continue
	}
	return connection, nil
}
