package main

import (
	"flag"
	"log"

	"github.com/laaksomavrick/microservices/pkg/amqp"
)

// Service should read varius user related events; do stuff
// Investigate request direct

var (
	uri          = flag.String("uri", "amqp://rabbitmq:rabbitmq@localhost:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queueName    = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("bindingKey", "test-key", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

func main() {
	c, err := amqp.NewConsumer(*uri, *exchangeName, *exchangeType, *queueName, *bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	forever := make(chan bool)

	// RPC (request/reply) is a popular pattern to implement with a messaging broker like RabbitMQ.
	// Tutorial 6 demonstrates its implementation with a variety of clients.
	// The typical way to do this is for RPC clients to send requests that are routed to a long lived (known) server queue.
	// The RPC server(s) consume requests from this queue and then send replies to each client using the queue named by the client in the reply-to header.
	go func() {
		for d := range c.Deliveries {
			log.Printf(
				"got %dB delivery: [%v] %s",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

	// log.Printf("shutting down")

	// if err := c.Shutdown(); err != nil {
	// 	log.Fatalf("error during shutdown: %s", err)
	// }

	// todo process.exit(1)
}
