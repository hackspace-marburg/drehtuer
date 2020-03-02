package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// mqttLogger wrapps mqtt's Logger interface to a wrapper for logrus.
type mqttLogger struct {
	ln func(v ...interface{})
	f  func(format string, v ...interface{})
}

func (l *mqttLogger) Println(v ...interface{}) {
	l.ln(v)
}

func (l *mqttLogger) Printf(format string, v ...interface{}) {
	l.f(format, v)
}

// newLogger creates a new mqttLogger, wrapping the two required methods.
func newLogger(ln func(v ...interface{}), f func(format string, v ...interface{})) *mqttLogger {
	return &mqttLogger{
		ln: ln,
		f:  f,
	}
}

// setupMqttLogger wraps mqtt's different Loggers to their logrus pedant.
func setupMqttLogger() {
	logger := log.StandardLogger()

	mqtt.CRITICAL = newLogger(logger.Panicln, logger.Panicf)
	mqtt.ERROR = newLogger(logger.Errorln, logger.Errorf)
	mqtt.WARN = newLogger(logger.Warnln, logger.Warnf)
	mqtt.DEBUG = newLogger(logger.Debugln, logger.Debugf)
}
