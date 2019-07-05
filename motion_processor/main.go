package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	//log filename and line :D
	// log.SetReportCaller(true)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

type MotionSensor interface {
	DetectMotion(done <-chan struct{}) <-chan bool
}

type MotionProcessor struct {
	// Period that the sensor remains active after detecting movement.
	SensorDetectionTimeout time.Duration
	LastEvent              *MotionEvent
}

func (m *MotionProcessor) ProcessStatusSequence(status bool) *MotionEvent {

	if status {
		m.LastEvent = &MotionEvent{
			Start: time.Now(),
		}
		return nil
	}

	if m.LastEvent == nil {
		return nil
	}

	event := m.LastEvent
	m.LastEvent = nil

	event.End = time.Now().Add(-m.SensorDetectionTimeout)
	event.Duration = event.End.Sub(event.Start)

	log.
		WithField("LastEventDuration", event.Duration.String()).
		Info("Event registered")

	return event
}

type MotionEvent struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

func main() {
	flag.UintVar(&sensorDetectionTimeout, "sensorDetectionTimeout", 0, "Time that the sensor remains active after detecting movement")

	flag.Parse()

	client := NewMqttClient()
	if token := client.Connect(); !token.WaitTimeout(time.Second*time.Duration(sendTimeout)) || token.Error() != nil {
		log.WithError(token.Error()).Panic("Unable to connect to mqtt server")
	}

	opts := client.OptionsReader()
	log.WithField("servers", opts.Servers()).Info("Connected to MQTT")

	processor := MotionProcessor{
		SensorDetectionTimeout: time.Second * time.Duration(sensorDetectionTimeout),
	}

	done := make(chan struct{})
	notification := DetectMotion(client, done)
	log.Info("Waiting for notifications")
	for {
		status := <-notification
		log.WithField("status", status).Info("New notificaton")
		event := processor.ProcessStatusSequence(status)
		if event != nil {
			jsonEvent, _ := json.Marshal(event)
			log.WithField("json", string(jsonEvent)).Debug("Json event generated")
			PublishMotionEvent(client, string(jsonEvent))
		}
	}
}
