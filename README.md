# golang

This repository contains the following go projects:

1. Power plant simulator
  - Distributed system that simulates data flow for a power plant.
  - It uses RabbitMQ to transmit messages from sensors to different receivers
  - Helpfull commands:
    - To start RabbitMQ server:
      - rabbitmq-server
    - To start sensor run:
      - go run src/powerplant/sensors/sensor.go -name sensor1 -max 9 -min 4
    - To start dispatcher run:
      - go run src/powerplant/dispatcher/main/main.go
    - To start webapp
      - go build src/powerplant/web/main.go
      - ./main