package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var prometheusListener string

var (
	promGaugeOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "drehtuer_door_open",
		Help: "Binary representation if the door is opened; or -1 if unknown",
	})

	promGaugeFlti = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "drehtuer_door_flti",
		Help: "Binary representation if currently are FLTI-only times; or -1 if unknown",
	})

	promGaugeLastUpdate = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "drehtuer_last_update_unix",
		Help: "Last update as the Unix timestamp in UTC; or -1 if unknown",
	})
)

// boolToFloat maps true to 1.0 and false to 0.0.
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	} else {
		return 0.0
	}
}

// setupPrometheus by starting a web server and exposing metrics.
func setupPrometheus() {
	if prometheusListener == "" {
		log.Debug("Skipping Prometheus setup as not configured")
		return
	}

	promGaugeOpen.Set(-1.0)
	promGaugeFlti.Set(-1.0)
	promGaugeLastUpdate.Set(-1.0)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(prometheusListener, nil); err != nil {
			log.Fatalf("Cannot start web server for Prometheus exporter, %v", err)
		}
	}()
}

// UpdatePromMetrics to be exposed through the web server.
func (doorState DoorState) UpdatePromMetrics() error {
	promGaugeOpen.Set(boolToFloat(doorState.Open))
	promGaugeFlti.Set(boolToFloat(doorState.FltiOnly))
	promGaugeLastUpdate.Set(float64(doorState.Timestamp))

	return nil
}
