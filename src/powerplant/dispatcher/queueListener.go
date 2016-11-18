package dispatcher

import (
	"log"
	"powerplant/dto"
	"powerplant/queueutils"

	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

const serverURL = "amqp://guest@localhost:5672"

type QueueListener struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	sources    map[string]<-chan amqp.Delivery
	dispatcher *EventDispatcher
}

func NewQueueListener(dispatcher *EventDispatcher) *QueueListener {
	qListener := QueueListener{}
	qListener.dispatcher = dispatcher
	qListener.sources = make(map[string]<-chan amqp.Delivery)
	qListener.connection, qListener.channel = queueutils.GetChannel(serverURL)
	return &qListener
}

// Sends a discovery call to signal all sensors to start sending data
func (ql *QueueListener) DiscoverSensors() {
	ql.channel.ExchangeDeclare(queueutils.SensorDiscoveryExchange, "fanout", false, false, false, false, nil)

	ql.channel.Publish(queueutils.SensorDiscoveryExchange, "", false, false, amqp.Publishing{})
}

func (ql *QueueListener) ListenForNewSource() {
	// create new temp queue to get fanout messages
	q := queueutils.GetQueue("", ql.channel, true)

	// bind queue to fanout exchange
	ql.channel.QueueBind(q.Name, "", "amq.fanout", false, nil)

	// read all messages from fanout exchange
	messages, _ := ql.channel.Consume(q.Name, "", true, false, false, false, nil)

	// send discovery call to find sensors
	ql.DiscoverSensors()

	// for every message create a new channel delivery and store it
	for message := range messages {
		qName := string(message.Body)
		ql.dispatcher.PublishEvent("DataSourceDiscovered", qName)
		if ql.sources[qName] == nil {
			sourceDelivery, _ := ql.channel.Consume(qName, "", true, false, false, false, nil)
			ql.sources[qName] = sourceDelivery

			go ql.addListener(sourceDelivery, qName)
		}
	}
}

func (ql *QueueListener) addListener(messages <-chan amqp.Delivery, qName string) {
	for message := range messages {
		reader := bytes.NewReader(message.Body)
		decoder := gob.NewDecoder(reader)
		sensorDto := new(dto.SensorMessage)
		decoder.Decode(sensorDto)

		log.Printf("Recieved message from %v: %v\n", qName, sensorDto)
		ql.dispatcher.PublishEvent("MessageReceived_"+message.RoutingKey, *sensorDto)
	}
}
