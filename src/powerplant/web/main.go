package main

import (
	"log"
	"net/http"
	"powerplant/web/controller"
)

func main() {
	controller.Initialize()
	log.Println("Server started on port: 3000. Navigate to http://localhost:3000/public")
	http.ListenAndServe(":3000", nil)
}
