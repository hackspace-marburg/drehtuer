package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var mqttClient mqtt.Client

// setupMqtt starts the MQTT client and subscribes for door.
func setupMqtt() {
	mqttOpts := mqtt.NewClientOptions().
		AddBroker("tcp://b2s.hsmr.cc:1883").
		SetClientID("drehtuer").
		SetKeepAlive(5 * time.Second)

	mqttClient = mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := mqttClient.Subscribe("door", 0, handleDoorMessage); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// teardownMqtt closes the MQTT client.
func teardownMqtt() {
	if token := mqttClient.Unsubscribe("door"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	mqttClient.Disconnect(1000)
}
