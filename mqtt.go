package main

import (
	"crypto/rand"
	"time"

	"github.com/akamensky/base58"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var mqttBroker string
var mqttClient mqtt.Client

// setupMqtt starts the MQTT client and subscribes for door.
func setupMqtt() {
	clientIdBuff := make([]byte, 6)
	if _, err := rand.Read(clientIdBuff); err != nil {
		log.Fatalf("Cannot generate a random ID, %v", err)
	}

	mqttOpts := mqtt.NewClientOptions().
		AddBroker(mqttBroker).
		SetClientID("drehtuer-" + string(base58.Encode(clientIdBuff))).
		SetKeepAlive(5 * time.Second)

	mqttClient = mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Cannot connect to MQTT, %v", token.Error())
	}

	if token := mqttClient.Subscribe("door", 0, handleDoorMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("Cannot subscribe to MQTT topic, %v", token.Error())
	}
}

// teardownMqtt closes the MQTT client.
func teardownMqtt() {
	if token := mqttClient.Unsubscribe("door"); token.Wait() && token.Error() != nil {
		log.Fatalf("Cannot unsubscribe from MQTT topic, %v", token.Error())
	}

	mqttClient.Disconnect(1000)
}
