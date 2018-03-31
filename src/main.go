package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go client()
	go server()

	var a string
	fmt.Scanln(&a)
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue,
		"",     // consumer,
		true,   // autoAck,
		false,  // exclusive,
		false,  //,noLocal,
		false,  //noWait,
		nil)
	failOnError(err, "Failed to register a consume")
	for msg := range msgs {
		log.Printf("Received message with: %s", msg.Body)
	}
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}
	for {
		ch.Publish("", q.Name, false, false, msg)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare("hello",
		false, //durable bool, //should i save to file system?
		false, //autoDelete bool, //no active consumers? should i delete
		false, //exclusive bool, //only accessible to the connection that requests it.
		false, //noWait bool, //only return existing queue
		nil)   //args amqp.Table)
	failOnError(err, "")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
		panic(fmt.Sprintf("%s : %s", msg, err))
	}
}
