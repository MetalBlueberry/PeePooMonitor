package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/metalblueberry/PeePooMonitor/motion_processor/mqtt"
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

type MotionEvent struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

func main() {
	hostname, _ := os.Hostname()
	server := flag.String("server", "tcp://mosquitto:1883", "The full URL of the MQTT server to connect to")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "guest", "A username to authenticate to the MQTT server")
	password := flag.String("password", "guest", "Password to match username")
	sendTimeout := flag.Uint("sendTimeout", 1, "Seconds to wait before failing to send a message")
	sensorDetectionTimeout := flag.Uint("sensorDetectionTimeout", 55, "Time that the sensor remains active after detecting movement")

	flag.Parse()

	clientOptions := &mqtt.MqttClientOptions{
		Server:      *server,
		Clientid:    *clientid,
		Username:    *username,
		Password:    *password,
		SendTimeout: *sendTimeout,
		OnConnect: func(client *mqtt.MqttClient) {
			log.Info("Reconnected")
		},
	}
	client := mqtt.NewMqttClient(clientOptions)
	err := client.Connect()
	if err != nil {
		log.WithError(err).Panic("Unable to connect to mqtt server")
	}

	processor := MotionProcessor{
		SensorDetectionTimeout: time.Second * time.Duration(*sensorDetectionTimeout),
	}

	done := make(chan struct{})
	notification := client.DetectMotion(done)
	for {
		status := <-notification
		log.WithField("status", status).Info("New notificaton")
		switch {
		case status:
			processor.LastEvent = &MotionEvent{
				Start: time.Now(),
			}
		case !status:
			if processor.LastEvent == nil {
				continue
			}
			processor.LastEvent.End = time.Now().Add(-processor.SensorDetectionTimeout)
			processor.LastEvent.Duration = processor.LastEvent.End.Sub(processor.LastEvent.Start)

			log.
				WithField("LastEventDuration", processor.LastEvent.Duration.String()).
				Info("Event registered")

			jsonEvent, _ := json.Marshal(processor.LastEvent)

			log.WithField("json", string(jsonEvent)).Debug("Json event generated")

			client.PublishMotionEvent(string(jsonEvent))
		}
	}

}
