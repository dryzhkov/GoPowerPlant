package dispatcher

import (
	"bytes"
	"encoding/gob"
	"powerplant/dto"
	"powerplant/queueutils"

	"github.com/streadway/amqp"
)

type WebappConnector struct {
	dispatcher IDispatcher
	connection *amqp.Connection
	channel    *amqp.Channel
	sources    []string
}

func NewWebappConnector(d IDispatcher) *WebappConnector {
	conn, ch := queueutils.GetChannel(queueutils.ServerURL)
	queueutils.GetQueue(queueutils.PersistDataQueue, ch, false)

	wc := WebappConnector{
		dispatcher: d,
		connection: conn,
		channel:    ch,
	}

	go wc.ListenForDiscoveryRequests()

	wc.dispatcher.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		wc.SubscribeToDataEvents(eventData.(string))
	})

	//init exchanges
	wc.channel.ExchangeDeclare(queueutils.WebappSourceExchange, "fanout", false, false, false, false, nil)
	wc.channel.ExchangeDeclare(queueutils.WebappReadingsExchange, "fanout", false, false, false, false, nil)
	return &wc
}

func (wc *WebappConnector) ListenForDiscoveryRequests() {
	q := queueutils.GetQueue(queueutils.WebappDiscoveryQueue, wc.channel, false)

	messages, _ := wc.channel.Consume(q.Name, "", true, false, false, false, nil)

	for range messages {
		for _, source := range wc.sources {
			wc.SendMessageSource(source)
		}
	}
}

func (wc *WebappConnector) SendMessageSource(source string) {
	wc.channel.Publish(queueutils.WebappSourceExchange,
		"",
		false,
		false,
		amqp.Publishing{
			Body: []byte(source),
		})
}

func (wc *WebappConnector) SubscribeToDataEvents(data string) {
	// return if we already know about the source
	for _, source := range wc.sources {
		if source == data {
			return
		}
	}

	wc.sources = append(wc.sources, data)
	wc.SendMessageSource(data)

	wc.dispatcher.AddListener("MessageReceived_"+data, func() func(interface{}) {
		return func(eventData interface{}) {
			sensorMessage := eventData.(dto.SensorMessage)
			buffer := new(bytes.Buffer)
			encoding := gob.NewEncoder(buffer)
			encoding.Encode(sensorMessage)
			msg := amqp.Publishing{
				Body: buffer.Bytes(),
			}
			wc.channel.Publish(queueutils.WebappReadingsExchange, "", false, false, msg)
		}
	}())
}
