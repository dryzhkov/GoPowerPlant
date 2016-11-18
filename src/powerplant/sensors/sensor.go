package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"powerplant/dto"
	"powerplant/queueutils"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var name = flag.String("name", "sensor", "name of the sensor")
var frequency = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5., "max value the sensor can generate")
var min = flag.Float64("min", 1., "min value the sensor can generate")
var step = flag.Float64("step", 0.1, "change allowed per measurement")

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var value float64
var nominal float64

func main() {
	flag.Parse()

	value = rand.Float64()*(*max-*min) + *min
	nominal = (*max-*min)/2 + *min

	conn, channel := queueutils.GetChannel(queueutils.ServerURL)
	defer conn.Close()
	defer channel.Close()

	messageQueue := queueutils.GetQueue(*name, channel, false)
	publishQueueName(channel)

	log.Printf("Listening for discovery calls.\n")
	go listenForDiscoveryCall(channel)

	ratio := 1000 / int(*frequency)
	duration, _ := time.ParseDuration(strconv.Itoa(ratio) + "ms")
	signal := time.Tick(duration)
	buffer := new(bytes.Buffer)
	encoding := gob.NewEncoder(buffer)

	for range signal {
		calcValue()
		message := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buffer.Reset()
		encoding = gob.NewEncoder(buffer)
		encoding.Encode(message)

		msg := amqp.Publishing{
			Body: buffer.Bytes(),
		}

		channel.Publish(
			"",                //exchange
			messageQueue.Name, //key
			false,             //mandatory
			false,             //immediate
			msg)               //message

		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func calcValue() {
	var maxStep, minStep float64

	if value < nominal {
		maxStep = *step
		minStep = -1 * *step * (value - *min) / (nominal - *min)
	} else {
		maxStep = *step * (*max - value) / (*max - nominal)
		minStep = -1 * *step
	}

	value += random.Float64()*(maxStep-minStep) + minStep
}

func publishQueueName(channel *amqp.Channel) {
	msg := amqp.Publishing{
		Body: []byte(*name),
	}

	channel.Publish(
		"amq.fanout", //exchange
		"",           //key
		false,        //mandatory
		false,        //immediate
		msg)          //message
}

func listenForDiscoveryCall(channel *amqp.Channel) {
	// create new temp queue to get fanout messages
	q := queueutils.GetQueue("", channel, true)

	// bind queue to fanout exchange
	channel.QueueBind(q.Name, "", queueutils.SensorDiscoveryExchange, false, nil)

	// read all messages from fanout exchange
	messages, _ := channel.Consume(q.Name, "", true, false, false, false, nil)

	for range messages {
		publishQueueName(channel)
	}
}
