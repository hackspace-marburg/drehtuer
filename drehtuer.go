package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// DoorState represents the door state's JSON as sent over the MQTT.
type DoorState struct {
	Timestamp int64 `json:"timestamp"`
	FltiOnly  bool  `json:"flti_only"`
	Open      bool  `json:"door_open"`
}

func (doorState DoorState) String() string {
	return fmt.Sprintf("Time: %v, FLTI only: %t, Open: %t",
		time.Unix(doorState.Timestamp, 0), doorState.FltiOnly, doorState.Open)
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

		if err := doorState.Publish(); err != nil {
			log.WithError(err).Error("Publishing to InfluxDB errored")
		}
	}
}

// waitSigint blocks the current thread until a SIGINT appears.
func waitSigint() {
	signalSyn := make(chan os.Signal)

	signal.Notify(signalSyn, os.Interrupt)
	<-signalSyn

	log.Info("Received SIGINT, closing down..")
}

func init() {
	var debugFlag bool

	flag.BoolVar(&debugFlag, "verbose", false, "Verbose logging output")
	flag.StringVar(&influxAddr, "influx", "", "InfluxDB address")
	flag.StringVar(&mqttBroker, "mqtt", "", "MQTT broker")

	flag.Parse()

	if debugFlag {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}

	if influxAddr == "" || mqttBroker == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	setupMqttLogger()
	setupMqtt()

	waitSigint()

	teardownMqtt()
}
