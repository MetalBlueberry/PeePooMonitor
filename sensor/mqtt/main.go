package mqtt

import (
	"crypto/tls"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type publisher interface {
	Connect() MQTT.Token
	Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token
	Disconnect(quiesce uint)
}
type Publisher interface {
	Disconnect(quiesce uint)
	PublishSensorStatus(status string)
	PublishPowerStatus(status bool)
}

type MqttClient struct {
	Server       string
	Topic        string
	Qos          int
	Clientid     string
	Username     string
	Password     string
	SendTimeout  uint
	CleanSession bool
	OnConnect    func(Publisher)
	client       publisher
}

func (m *MqttClient) Connect() {
	connOpts := MQTT.
		NewClientOptions().
		AddBroker(m.Server).
		SetClientID(m.Clientid).
		SetCleanSession(m.CleanSession).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(client MQTT.Client) {
			m.OnConnect(m)
		})

	if m.Username != "" {
		connOpts.SetUsername(m.Username)
		if m.Password != "" {
			connOpts.SetPassword(m.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	m.client = MQTT.NewClient(connOpts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
		return
	}
	log.WithField("Server", m.Server).Info("Connected to Server")
}

func (m *MqttClient) Publish(message string, topic string, retained bool) bool {
	return m.client.Publish(topic, byte(m.Qos), retained, message).WaitTimeout(time.Second * time.Duration(m.SendTimeout))
}

func (m *MqttClient) Disconnect(quiesce uint) {
	m.client.Disconnect(100)
}

func (m *MqttClient) PublishSensorStatus(status string) {
	m.Publish(status, "/devices/sensor/raw/status", true)
}
func (m *MqttClient) PublishPowerStatus(status bool) {
	if status {
		m.Publish("PowerOn", "/devices/sensor/connection/status", true)
	} else {
		m.Publish("PowerOff", "/devices/sensor/connection/status", true)
	}
}
