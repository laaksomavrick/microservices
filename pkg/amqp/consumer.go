package amqp

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	Done       chan error
	Deliveries <-chan amqp.Delivery
	tag        string
}

func (consumer *Consumer) Hello() {
	fmt.Print("hello from consumer")
}

func NewConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		Conn:       nil,
		Channel:    nil,
		Done:       make(chan error),
		Deliveries: make(<-chan amqp.Delivery),
		tag:        ctag,
	}

	var err error

	log.Printf("dialing %s", amqpURI)
	c.Conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	log.Printf("got Connection, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%s)", exchange)
	if err = c.Channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue (%s)", queue)
	state, err := c.Channel.QueueDeclare(
		queue, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		state.Messages, state.Consumers, key)

	if err = c.Channel.QueueBind(
		queue,    // name of the queue
		key,      // bindingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag '%s')", c.tag)
	c.Deliveries, err = c.Channel.Consume(
		queue, // name
		c.tag, // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	// go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handler to exit
	return <-c.Done
}
