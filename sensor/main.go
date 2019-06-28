package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/metalblueberry/PeePooMonitor/sensor/hcsr51"
	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/host/bcm283x"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	//log filename and line :D
	log.SetReportCaller(true)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

var sensor *hcsr51.HCSR51

func main() {

	flag.Parse()

	sensor := hcsr51.NewHCSR51(bcm283x.GPIO17)
	client := NewMqttClient(sensor)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Panic("Unable to connect to mqtt server")
	}
	opts := client.OptionsReader()
	log.WithField("servers", opts.Servers()).Info("Connected to MQTT")

	done := make(chan struct{})
	defer close(done)
	notifier := sensor.DetectMotion(done)
	PublishPowerStatus(client, true)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		signal.Reset(os.Interrupt, syscall.SIGTERM)
		log.Info("Trying to shutdown the system, press ctrl+c again to force stop")
		done <- struct{}{}
	}()

	for {
		status, ok := <-notifier
		if !ok {
			log.Info("notification channel is closed")
			break
		}

		if err := PublishSensorStatus(client, status.String()); err != nil {
			log.WithError(err).Error("Unable to publish")
		}
		log.WithField("status", status).Info("Published status")
	}

	log.Trace("Notify PowerStatus")
	if err := PublishPowerStatus(client, false); err != nil {
		log.WithError(err).Error("Unable to publish")
	}
	log.Trace("Disconnect from mosquitto")
	client.Disconnect(100)
	log.Trace("Disconnect process finished")
}
