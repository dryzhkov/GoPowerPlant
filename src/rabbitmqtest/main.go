package main

import (
	"fmt"
	"log"
	"powerplant/queueutils"

	"github.com/streadway/amqp"
)

func main() {
	go server()
	go client()

	var a string
	fmt.Scanln(&a)
}

func client() {
	conn, channel, queue := getQueue()
	defer conn.Close()
	defer channel.Close()

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register client")

	for message := range messages {
		log.Printf("Recieved message: %s", message.Body)
	}
}

func server() {
	conn, channel, queue := getQueue()
	defer conn.Close()
	defer channel.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}
	for {
		channel.Publish(
			"",         //exchange
			queue.Name, //key
			false,      //mandatory
			false,      //immediate
			msg)        //message
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial(queueutils.ServerURL)
	failOnError(err, "Failed to open connection to RabbitMQ")
	channel, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	queue, err := channel.QueueDeclare(
		"hello", //name
		false,   //durable
		false,   //autoDelete
		false,   //exclusive
		false,   //noWait
		nil)     //args table
	failOnError(err, "Failed to connect to queue")

	return conn, channel, &queue
}

func failOnError(err error, message string) {
	if err != nil {
		// print error and stop application
		log.Fatalf("%s: %s", message, err)
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}
