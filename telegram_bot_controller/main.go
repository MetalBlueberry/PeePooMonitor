package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/metalblueberry/PeePooMonitor/telegram_bot_controller/mqtt"
	"github.com/metalblueberry/PeePooMonitor/telegram_bot_controller/tgbot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("hello")

	server := flag.String("server", "tcp://mosquitto:1883", "The full URL of the MQTT server to connect to")
	clientid := flag.String("clientid", "telegram_bot_controller", "A clientid for the connection")
	username := flag.String("username", "guest", "A username to authenticate to the MQTT server")
	password := flag.String("password", "guest", "Password to match username")
	sendTimeout := flag.Uint("sendTimeout", 1, "Seconds to wait before failing to send a message")

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

	TelegramToken, err := ioutil.ReadFile("/run/secrets/telegram_bot_token")
	if err != nil {
		log.WithError(err).Panic("cannot load telegram bot token from secrets")
	}

	bot := tgbot.TelegramBot{
		Token: strings.Trim(string(TelegramToken), " \n"),
	}
	done := make(chan struct{})
	notifier := bot.NotifyMotion()
	events := client.SubscribeToMotionEvents(done)

	for {
		event := <-events
		notifier <- fmt.Sprintf("Cat has been doing its things\nStart: %s\nEnd: %s\n%f seconds", event.Start.Format("2006-01-02 15:04:05"), event.End.Format("2006-01-02 15:04:05"), event.Duration.Seconds())
	}
}
