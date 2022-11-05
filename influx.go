package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var influxAddr string

// PublishInflux this DoorState to the InfluxDB.
func (doorState DoorState) PublishInflux() error {
	if influxAddr == "" {
		log.Debug("Skipping publishing InfluxDB as not configured")
		return nil
	}

	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: influxAddr,
	})

	if err != nil {
		return err
	}
	defer client.Close()

	state := map[string]interface{}{
		"flti": doorState.FltiOnly,
		"open": doorState.Open,
	}

	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: "doorstate",
	})
	pt, err := influx.NewPoint("door", nil, state, time.Unix(doorState.Timestamp, 0))
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	return client.Write(bp)
}
