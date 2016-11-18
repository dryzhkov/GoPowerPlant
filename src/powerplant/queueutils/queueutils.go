package queueutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const ServerURL = "amqp://guest@localhost:5672"

const SensorDiscoveryExchange string = "SensorDiscovery"
const PersistDataQueue string = "PersistData"
const WebappSourceExchange string = "WebappSources"
const WebappReadingsExchange string = "WebappReadings"
const WebappDiscoveryQueue string = "WebappDiscovery"

func failOnError(err error, message string) {
	if err != nil {
		// print error and stop application
		log.Fatalf("%s: %s", message, err)
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}

// GetChannel - takes URL of RabbitMQ server, opens a connection and returns connection and channel reference
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to open connection to RabbitMQ")
	channel, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	return conn, channel
}

// GetQueue - given a name and channel reference, creates and returns a new queueutils
func GetQueue(name string, channel *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := channel.QueueDeclare(
		name,       //name
		false,      //durable
		autoDelete, //autoDelete
		false,      //exclusive
		false,      //noWait
		nil)        //args table
	failOnError(err, "Failed to create a queue")
	return &q
}
