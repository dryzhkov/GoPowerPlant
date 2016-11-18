package main

import (
	"fmt"
	"powerplant/dispatcher"
)

var wc *dispatcher.WebappConnector

func main() {
	d := dispatcher.NewEventDispatcher()
	wc = dispatcher.NewWebappConnector(d)
	qListener := dispatcher.NewQueueListener(d)
	go qListener.ListenForNewSource()

	var x string
	fmt.Scanln(&x)
}
