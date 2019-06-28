package main

import (
	"crypto/tls"
	"flag"
	"os"
	"time"

	"github.com/metalblueberry/PeePooMonitor/sensor/hcsr51"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Addresses string

const (
	AddressSensorStatus Addresses = "/devices/sensor/raw/status"
	AddressPowerStatus  Addresses = "/devices/sensor/connection/status"
)

var (
	server      string
	clientid    string
	username    string
	password    string
	sendTimeout uint
)

func init() {
	hostname, _ := os.Hostname()
	flag.StringVar(&server, "server", "tcp://mosquitto:1883", "The full URL of the MQTT server to connect to")
	flag.StringVar(&clientid, "clientid", hostname, "A clientid for the connection")
	flag.StringVar(&username, "username", "guest", "A username to authenticate to the MQTT server")
	flag.StringVar(&password, "password", "guest", "Password to match username")
	flag.UintVar(&sendTimeout, "sendTimeout", 5, "Seconds to wait before failing to send a message")
}

func NewMqttClient(sensor *hcsr51.HCSR51) MQTT.Client {
	connOpts := MQTT.
		NewClientOptions().
		AddBroker(server).
		SetClientID(clientid).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(client MQTT.Client) {
			log.Info("Connected")
			PublishPowerStatus(client, true)
			PublishSensorStatus(client, sensor.Status().String())
		}).
		SetConnectionLostHandler(func(client MQTT.Client, err error) {
			log.WithError(err).Warning("Disconnected from MQTT")
		})

	if username != "" {
		connOpts.SetUsername(username)
		if password != "" {
			connOpts.SetPassword(password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	return MQTT.NewClient(connOpts)
}

func PublishSensorStatus(client MQTT.Client, status string) error {
	token := client.Publish(string(AddressSensorStatus), 1, false, status)
	if !token.WaitTimeout(time.Second * time.Duration(sendTimeout)) {
		return errors.New("Timeout publishing message")
	}
	if err := token.Error(); err != nil {
		return errors.Wrap(err, "Error Publishing Sensor Status")
	}
	return nil
}
func PublishPowerStatus(client MQTT.Client, status bool) error {
	if status {
		token := client.Publish(string(AddressPowerStatus), 1, true, "PowerOn")
		if !token.WaitTimeout(time.Second * time.Duration(sendTimeout)) {
			return errors.New("Timeout publishing message")
		}
		if err := token.Error(); err != nil {
			return errors.Wrap(err, "Error Publishing Sensor Power Status")
		}
	} else {
		token := client.Publish(string(AddressPowerStatus), 1, true, "PowerOff")
		if !token.WaitTimeout(time.Second * time.Duration(sendTimeout)) {
			return errors.New("Timeout publishing message")
		}
		if err := token.Error(); err != nil {
			return errors.Wrap(err, "Error Publishing Sensor Power Status")
		}
	}
	return nil
}
