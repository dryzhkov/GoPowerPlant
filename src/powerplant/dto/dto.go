package dto

import "time"
import "encoding/gob"

// SensorMessage contains properties that are sent to RabbitMQ
type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
