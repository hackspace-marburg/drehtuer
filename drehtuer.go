package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// DoorState represents the door state's JSON as sent over the MQTT.
type DoorState struct {
	Timestamp uint64 `json:"timestamp"`
	FltiOnly  bool   `json:"flti_only"`
	Open      bool   `json:"door_open"`
}

func (doorState DoorState) String() string {
	return fmt.Sprintf("Timestamp: %d, FLTI only: %t, Open: %t",
		doorState.Timestamp, doorState.FltiOnly, doorState.Open)
}

// handleDoorMessage is called from the MQTT client for new messages with the door topic.
func handleDoorMessage(_ mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received MQTT message")

	var doorState DoorState
	if err := json.Unmarshal(msg.Payload(), &doorState); err != nil {
		log.WithError(err).Error("Unmarshaling JSON errored")
	} else {
		log.WithField("door", doorState).Info("Received MQTT door state")
	}
}

// waitSigint blocks the current thread until a SIGINT appears.
func waitSigint() {
	signalSyn := make(chan os.Signal)

	signal.Notify(signalSyn, os.Interrupt)
	<-signalSyn

	log.Info("Received SIGINT, closing down..")
}

func main() {
	setupMqttLogger()
	setupMqtt()

	waitSigint()

	teardownMqtt()
}
