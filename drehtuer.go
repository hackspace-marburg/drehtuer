package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func handleDoorMessage(_ mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Info("Received MQTT message")
}

func main() {
	setupMqttLogger()

	mqttOpts := mqtt.NewClientOptions().
		AddBroker("tcp://b2s.hsmr.cc:1883").
		SetClientID("drehtuer").
		SetKeepAlive(5 * time.Second)

	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := mqttClient.Subscribe("door", 0, handleDoorMessage); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	time.Sleep(10 * time.Second)

	if token := mqttClient.Unsubscribe("door"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	mqttClient.Disconnect(1000)
}
