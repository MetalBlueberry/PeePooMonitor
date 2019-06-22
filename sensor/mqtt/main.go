package mqtt

import (
	"crypto/tls"
	"time"

	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//go:generate mockgen -destination=../mocks/mock_$GOFILE -package=mocks  --source=$GOFILE

type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}
type mqttPublisher interface {
	Connect() MQTT.Token
	Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token
	Disconnect(quiesce uint)
}

// Publisher is an interface to report the status of the sensor.
type Publisher interface {
	Disconnect(quiesce uint)
	PublishSensorStatus(status string)
	PublishPowerStatus(status bool)
}

type MqttClientOptions struct {
	Server       string
	Qos          int
	Clientid     string
	Username     string
	Password     string
	SendTimeout  uint
	CleanSession bool
	OnConnect    func(Publisher)
}

type MqttClient struct {
	Client  mqttPublisher
	Options *MqttClientOptions
}

type Addresses string

const (
	SensorStatusAddress Addresses = "/devices/sensor/raw/status"
	PowerStatusAddress  Addresses = "/devices/sensor/connection/status"
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

func (m *MqttClient) Publish(message string, topic string, retained bool) bool {
	log.
		WithField("topic", topic).
		WithField("message", message).
		Debug("Publishing message")

	return m.Client.Publish(topic, byte(m.Options.Qos), retained, message).WaitTimeout(time.Second * time.Duration(m.Options.SendTimeout))
}

func (m *MqttClient) Disconnect(quiesce uint) {
	m.Client.Disconnect(100)
}

func (m *MqttClient) PublishSensorStatus(status string) {
	m.Publish(status, string(SensorStatusAddress), true)
}
func (m *MqttClient) PublishPowerStatus(status bool) {
	if status {
		m.Publish("PowerOn", string(PowerStatusAddress), true)
	} else {
		m.Publish("PowerOff", string(PowerStatusAddress), true)
	}
}
