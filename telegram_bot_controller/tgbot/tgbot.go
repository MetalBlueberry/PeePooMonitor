package tgbot

import (

	//"github.com/metalblueberry/PeePooMonitor/pkg/telegrambot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type TelegramBot struct {
	Token string
	bot   *tgbotapi.BotAPI
}

var (
	admins = []int64{
		// 7452294,
		142825882,
	}
)

func (t TelegramBot) NotifyMotion() chan<- string {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}
	t.bot = bot

	log.WithField("UserName", bot.Self.UserName).Info("Authorized on account")

	notifications := make(chan string)

	go func(notifications <-chan string) {
		for {
			text := <-notifications
			for _, ChatID := range admins {
				msg := tgbotapi.NewMessage(ChatID, text)
				_, err := t.bot.Send(msg)

				if err != nil {
					log.WithError(err).Error("Error sending telegram message")
				}
			}
		}

	}(notifications)

	return notifications
}
