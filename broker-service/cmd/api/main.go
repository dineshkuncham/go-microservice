package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit *amqp.Connection
}

const webport = "8080"

func main() {
	connection, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	app := Config{
		Rabbit: connection,
	}

	log.Printf("Starting the broker service on port %s", webport)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
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
