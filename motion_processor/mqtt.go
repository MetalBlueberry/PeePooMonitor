package main

import (
	"crypto/tls"
	"flag"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Addresses string

const (
	AddressSensorStatus Addresses = "/devices/sensor/raw/status"
	AddressPowerStatus  Addresses = "/devices/sensor/connection/status"
	AddressMotionEvent  Addresses = "/events/motion"
)

var (
	server                 string
	clientid               string
	username               string
	password               string
	sendTimeout            uint
	sensorDetectionTimeout uint
)

func init() {
	flag.StringVar(&server, "server", "tcp://mosquitto:1883", "The full URL of the MQTT server to connect to")
	flag.StringVar(&clientid, "clientid", "motion_processor", "A clientid for the connection")
	flag.StringVar(&username, "username", "guest", "A username to authenticate to the MQTT server")
	flag.StringVar(&password, "password", "guest", "Password to match username")
	flag.UintVar(&sendTimeout, "sendTimeout", 5, "Seconds to wait before failing to send a message")
}

func NewMqttClient() MQTT.Client {
	connOpts := MQTT.
		NewClientOptions().
		AddBroker(server).
		SetClientID(clientid).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(client MQTT.Client) {
			log.Info("Connected")
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

func DetectMotion(client MQTT.Client, done <-chan struct{}) <-chan bool {
	notify := make(chan bool)
	client.Subscribe(string(AddressSensorStatus), 1, GenerateMessageHandler(notify))
	go func(done <-chan struct{}) {
		<-done
		close(notify)
		token := client.Unsubscribe(string(AddressSensorStatus))
		if !token.WaitTimeout(time.Second * time.Duration(sendTimeout)) {
			log.Warn("Timeout unsubscribing")
		}
		if err := token.Error(); err != nil {
			log.WithError(err).Warn("Error when unsubscribing channel")
		}
	}(done)
	return notify
}

func GenerateMessageHandler(notify chan<- bool) MQTT.MessageHandler {
	return func(client MQTT.Client, message MQTT.Message) {
		log.
			WithField("payload", string(message.Payload())).
			WithField("topic", message.Topic()).
			Debug("New Message")
		notify <- string(message.Payload()) == "High"
	}
}

func PublishMotionEvent(client MQTT.Client, jsonEvent string) error {
	token := client.Publish(string(AddressMotionEvent), 1, true, "PowerOn")
	if !token.WaitTimeout(time.Second * time.Duration(sendTimeout)) {
		return errors.New("Timeout publishing message")
	}
	if err := token.Error(); err != nil {
		return errors.Wrap(err, "Error Publishing Sensor Power Status")
	}
	return nil
}
