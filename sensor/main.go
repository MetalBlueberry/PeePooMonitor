package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/metalblueberry/PeePooMonitor/sensor/hcsr51"
	"github.com/metalblueberry/PeePooMonitor/sensor/mqtt"
	"periph.io/x/periph/host/bcm283x"
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

func main() {

	hostname, _ := os.Hostname()
	server := flag.String("server", "tcp://mosquitto:1883", "The full URL of the MQTT server to connect to")
	qos := flag.Int("qos", 1, "The QoS to send the messages at")
	clientid := flag.String("clientid", hostname, "A clientid for the connection")
	username := flag.String("username", "guest", "A username to authenticate to the MQTT server")
	password := flag.String("password", "guest", "Password to match username")
	sendTimeout := flag.Uint("sendTimeout", 1, "Seconds to wait before failing to send a message")

	flag.Parse()

	sensor := hcsr51.NewHCSR51(bcm283x.GPIO17)
	clientOptions := &mqtt.MqttClientOptions{
		Server:      *server,
		Qos:         *qos,
		Clientid:    *clientid,
		Username:    *username,
		Password:    *password,
		SendTimeout: *sendTimeout,
		OnConnect: func(client mqtt.Publisher) {
			log.Info("Reconnected")
			status := sensor.Status()
			client.PublishPowerStatus(true)
			client.PublishSensorStatus(status.String())
		},
	}
	client := mqtt.NewMqttClient(clientOptions)
	client.Connect()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Debug("Disconnect")
		client.PublishPowerStatus(false)
		client.Disconnect(100)
		os.Exit(1)
	}()

	done := make(chan struct{})
	defer close(done)
	notifier := sensor.DetectMotion(done)

	client.Connect()
	client.PublishPowerStatus(true)

	for {
		status := <-notifier

		client.PublishSensorStatus(status.String())
		log.WithField("status", status).Info("Published status")
	}

}
