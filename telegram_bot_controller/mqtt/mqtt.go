package mqtt

import (
	"encoding/json"
	"crypto/tls"
	"time"

	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//go:generate mockgen -destination=../mocks/mock_$GOFILE -package=mocks  --source=$GOFILE

type token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

type MqttClientOptions struct {
	Server       string
	Clientid     string
	Username     string
	Password     string
	SendTimeout  uint
	CleanSession bool
	OnConnect    func(*MqttClient)
}

type MqttClient struct {
	Client  MQTT.Client
	Options *MqttClientOptions
}

type MotionEvent struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type Addresses string

const (
	PowerStatusAddress Addresses = "/devices/sensor/connection/status"
	AddressMotionEvent Addresses = "/events/motion"
)

func NewMqttClient(opt *MqttClientOptions) *MqttClient {
	var client *MqttClient
	connOpts := MQTT.
		NewClientOptions().
		AddBroker(opt.Server).
		SetClientID(opt.Clientid).
		SetCleanSession(opt.CleanSession).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(MQTT.Client) {
			opt.OnConnect(client)
		})

	if opt.Username != "" {
		connOpts.SetUsername(opt.Username)
		if opt.Password != "" {
			connOpts.SetPassword(opt.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	realClient := MQTT.NewClient(connOpts)
	client = &MqttClient{
		Options: opt,
		Client:  realClient,
	}
	return client
}

func (m *MqttClient) Connect() error {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.WithField("Server", m.Options.Server).Info("Connected to Server")
	return nil
}

func (m *MqttClient) Publish(message string, qos byte, topic string, retained bool) bool {
	log.
		WithField("topic", topic).
		WithField("message", message).
		Debug("Publishing message")

	return m.Client.Publish(topic, qos, retained, message).WaitTimeout(time.Second * time.Duration(m.Options.SendTimeout))
}

type MessageHandler func(*MqttClient, MQTT.Message)

func (m *MqttClient) Subscribe(topic Addresses, qos byte, callback MessageHandler) MQTT.Token {
	return m.Client.Subscribe(string(topic), qos, func(client MQTT.Client, message MQTT.Message) {
		callback(m, message)
	})
}

func (m *MqttClient) Disconnect(quiesce uint) {
	m.Client.Disconnect(100)
}

func (m *MqttClient) Unsubscribe(topics ...Addresses) bool {
	for topic := range topics {
		ok := m.Client.Unsubscribe(string(topic)).WaitTimeout(time.Second * time.Duration(m.Options.SendTimeout))
		if !ok {
			return false
		}
	}
	return true
}

func (m *MqttClient) SubscribeToMotionEvents(done <-chan struct{}) <-chan MotionEvent {
	notify := make(chan MotionEvent)
	m.Subscribe(AddressMotionEvent, 1, GenerateMotionEventMessageHandler(notify))
	go func(done <-chan struct{}) {
		<-done
		close(notify)
		ok := m.Unsubscribe(AddressMotionEvent)
		if !ok {
			log.Warn("Unsubscribe action failed")
		}
	}(done)
	return notify
}

func GenerateMotionEventMessageHandler(notify chan<- MotionEvent) MessageHandler {
	return func(client *MqttClient, message MQTT.Message) {
		log.
			WithField("payload", string(message.Payload())).
			WithField("topic", message.Topic()).
			Debug("New Message")

		var event MotionEvent 
		err := json.Unmarshal(message.Payload(), &event)
		if err != nil{
			log.
			WithError(err).
			WithField("payload", string(message.Payload())).
			WithField("topic", message.Topic()).
			Error("Unable to unmarshall the payload from motion event")
			return 
		}

		notify <- event
	}
}
